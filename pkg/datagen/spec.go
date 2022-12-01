package datagen

import (
	"encoding/json"
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/charset"
	"github.com/flipkart-incubator/diligent/pkg/keygen"
	"github.com/flipkart-incubator/diligent/pkg/strgen"
	"github.com/flipkart-incubator/diligent/pkg/strtr"
	"io/ioutil"
	"log"
	"os"
)

const (
	specTypeString = "diligent/schema-a"
	currentSpecVersion = 1
)

// Spec is the specification for generating a DataGen for the following schema:
//    pk        varchar,
//    uniq      varchar,
//    small_grp varchar,
//    large_grp varchar,
//    same      varchar,
//    seq_num   int    ,
//    ts        timestamp,
//    payload   varchar,
type Spec struct {
	SpecType       string
	Version        int
	RecordSize     int
	KeyGenSpec     *keygen.LeveledKeyGenSpec
	UniqTrSpec     *strtr.Spec
	SmallGrpTrSpec *strtr.Spec
	LargeGrpTrSpec *strtr.Spec
	FixedValue     string
	ReusePayload   bool
}

func NewSpec(recordCount int, recordSize int, reusePayload bool) *Spec {
	numSubKeysPerLevel := computeNumSubKeysPerLevel(recordCount)

	kgSpec := keygen.NewRandomLeveledKeyGenSpec(numSubKeysPerLevel, 5)
	strGen := strgen.NewStrGen(charset.AlphaUp)

	ds := &Spec{
		SpecType:       specTypeString,
		Version:        currentSpecVersion,
		RecordSize:     recordSize,
		KeyGenSpec:     kgSpec,
		UniqTrSpec:     strtr.NewRandomTrSpec(charset.AlphaUp),
		SmallGrpTrSpec: strtr.NewRandomTrSpec(charset.AlphaUp),
		LargeGrpTrSpec: strtr.NewRandomTrSpec(charset.AlphaUp),
		FixedValue:     strGen.RandomString(kgSpec.KeyLength()),
		ReusePayload:   reusePayload,
	}

	return ds
}

func LoadSpecFromFile(fileName string) (*Spec, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	spec := &Spec{}
	err = json.Unmarshal(bytes, spec)
	if err != nil {
		return nil, err
	}
	if !spec.IsValid() {
		return nil, fmt.Errorf("file did not contain a valid spec")
	}
	return spec, nil
}

func (spec *Spec) IsValid() bool {
	if spec.SpecType != specTypeString {
		return false
	}
	if spec.Version != 1 {
		return false
	}
	if spec.RecordSize <= 0 {
		return false
	}
	if !spec.KeyGenSpec.IsValid() {
		return false
	}
	if !spec.UniqTrSpec.IsValid() {
		return false
	}
	if !spec.SmallGrpTrSpec.IsValid() {
		return false
	}
	if !spec.LargeGrpTrSpec.IsValid() {
		return false
	}
	if len(spec.FixedValue) != spec.KeyGenSpec.KeyLength() {
		return false
	}
	return true
}

func (spec *Spec) SaveToFile(fileName string) error {
	jsonBytes := spec.Json()
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonBytes)
	if err != nil {
		return err
	}
	return nil
}

func (spec *Spec) Json() []byte {
	bytes, err := json.MarshalIndent(spec, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

// We create a LeveledKeyGen with 3 levels
// Generates the number of subkeys in each level based on the total number of records
func computeNumSubKeysPerLevel(recordCount int) []int {
	var subKeysPerLevel []int

	switch {
	case recordCount < 1:
		panic(fmt.Sprintf("Invalid record count: %d", recordCount))
	case recordCount < 10:
		// Single digit record count
		subKeysPerLevel = []int{1, 1, recordCount}
	case recordCount < 100:
		// Two digit record count rounded to 10s
		recordCount = recordCount - (recordCount % 10)
		subKeysPerLevel = []int{1, recordCount / 10, 10}
	case recordCount < 1000:
		// Three digit record count rounded to 100s
		recordCount = recordCount - (recordCount % 100)
		subKeysPerLevel = []int{recordCount / 100, 10, 10}
	case recordCount < 10_000:
		recordCount = recordCount - (recordCount % 1000)
		subKeysPerLevel = []int{recordCount / 1000, 100, 10}
	default:
		// 10K or more records rounded to 10Ks
		recordCount = recordCount - (recordCount % 10000)
		subKeysPerLevel = []int{recordCount / 10_000, 1000, 10}
	}

	// Validate the output produced

	if len(subKeysPerLevel) < 1 {
		panic(fmt.Sprintf("Invalid subKeysPerLevel spec generated: %#v", subKeysPerLevel))
	}
	rc := 1
	for _, n := range subKeysPerLevel {
		if n < 1 {
			panic(fmt.Sprintf("Invalid subKeysPerLevel spec generated: %#v", subKeysPerLevel))
		}
		rc *= n
	}
	if rc != recordCount {
		panic(fmt.Sprintf("KeySet size mismatch. Expected %d, Got %d, numSubKeysPerLevel=%#v",
			recordCount, rc, subKeysPerLevel))
	}
	return subKeysPerLevel
}
