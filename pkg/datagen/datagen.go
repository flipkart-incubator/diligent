package datagen

import (
	"github.com/flipkart-incubator/diligent/pkg/charset"
	"github.com/flipkart-incubator/diligent/pkg/keygen"
	"github.com/flipkart-incubator/diligent/pkg/strgen"
	"github.com/flipkart-incubator/diligent/pkg/strtr"
	"go.uber.org/atomic"
	"math"
	"time"
)

// Record represents a single record with a fixed, well-defined set of fields
type Record struct {
	Pk         string
	Uniq       string
	SmallGrp   string
	LargeGrp   string
	FixedValue string
	SeqNum     int64
	TimeStamp  int64
	Payload    string
}

// DataGen can generate a bunch of data of type Record
type DataGen struct {
	recordSize  int
	payloadSize int
	keyGen      *keygen.LeveledKeyGen
	uniqTr      *strtr.Tr
	smallGrpTr  *strtr.Tr
	largeGrpTr  *strtr.Tr
	fixedValue  string
	seqGen      *atomic.Int64
	payloadGen  *strgen.StrGen
}

func NewDataGen(spec *Spec) *DataGen {
	if !spec.IsValid() {
		panic("Invalid spec for DataGen")
	}

	nonPayloadSize :=
		spec.KeyGenSpec.KeyLength()*2 /* pk, uniq */ +
			spec.KeyGenSpec.LengthOfKeyPrefixAtLevel(1) /* smallGrp */ +
			spec.KeyGenSpec.LengthOfKeyPrefixAtLevel(0) /* largeGrp */ +
			len(spec.FixedValue) /* fixed value */ +
			4 + 8 /* seq, ts */

	payloadSize := int(math.Max(float64(spec.RecordSize-nonPayloadSize), 1))

	return &DataGen{
		recordSize:  spec.RecordSize,
		payloadSize: payloadSize,
		keyGen:      keygen.NewLeveledKeyGen(spec.KeyGenSpec),
		uniqTr:      strtr.NewTr(spec.UniqTrSpec),
		smallGrpTr:  strtr.NewTr(spec.SmallGrpTrSpec),
		largeGrpTr:  strtr.NewTr(spec.LargeGrpTrSpec),
		fixedValue:  spec.FixedValue,
		seqGen:      atomic.NewInt64(0),
		payloadGen:  strgen.NewStrGen(charset.AlphaUp),
	}
}

func (dg *DataGen) NumRecords() int {
	return dg.keyGen.NumKeys()
}

func (dg *DataGen) Key(n int) string {
	return dg.keyGen.Key(n).String()
}

func (dg *DataGen) Uniq(n int) string {
	return dg.uniqTr.Apply(dg.keyGen.Key(n).String())
}

func (dg *DataGen) SmallGrp(n int) string {
	return dg.smallGrpTr.Apply(dg.keyGen.Key(n).Prefix(1))
}

func (dg *DataGen) LargeGrp(n int) string {
	return dg.largeGrpTr.Apply(dg.keyGen.Key(n).Prefix(0))
}

func (dg *DataGen) FixedValue() string {
	return dg.fixedValue
}

func (dg *DataGen) RandomPayload() string {
	return dg.payloadGen.RandomString(dg.payloadSize)
}

func (dg *DataGen) Record(n int) *Record {
	key := dg.keyGen.Key(n)
	return &Record{
		Pk:         key.String(),
		Uniq:       dg.uniqTr.Apply(key.String()),
		SmallGrp:   dg.smallGrpTr.Apply(key.Prefix(1)),
		LargeGrp:   dg.largeGrpTr.Apply(key.Prefix(0)),
		FixedValue: dg.fixedValue,
		SeqNum:     dg.seqGen.Inc(),
		TimeStamp:  time.Now().Unix(),
		Payload:    dg.payloadGen.RandomString(dg.payloadSize),
	}
}
