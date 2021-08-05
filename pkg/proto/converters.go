package proto

import (
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/keygen"
	"github.com/flipkart-incubator/diligent/pkg/strtr"
)

// TODO: Add a hash check for DataSpec transfer
func DataSpecToProto(spec *datagen.Spec) *DataSpec {
	return &DataSpec{
		SpecType:       spec.SpecType,
		Version:        int32(spec.Version),
		RecordSize:     int32(spec.RecordSize),
		KeyGenSpec:     KeyGenSpecToProto(spec.KeyGenSpec),
		UniqTrSpec:     TrSpecToProto(spec.UniqTrSpec),
		SmallGrpTrSpec: TrSpecToProto(spec.SmallGrpTrSpec),
		LargeGrpTrSpec: TrSpecToProto(spec.LargeGrpTrSpec),
		FixedValue:     spec.FixedValue,
	}
}

func DataSpecFromProto(protoSpec *DataSpec) *datagen.Spec {
	return &datagen.Spec{
		SpecType:       protoSpec.GetSpecType(),
		Version:        int(protoSpec.GetVersion()),
		RecordSize:     int(protoSpec.GetRecordSize()),
		KeyGenSpec:     KeyGenSpecFromProto(protoSpec.GetKeyGenSpec()),
		UniqTrSpec:     TrSpecFromProto(protoSpec.GetUniqTrSpec()),
		SmallGrpTrSpec: TrSpecFromProto(protoSpec.GetSmallGrpTrSpec()),
		LargeGrpTrSpec: TrSpecFromProto(protoSpec.GetLargeGrpTrSpec()),
		FixedValue:     protoSpec.GetFixedValue(),
	}
}

func KeyGenSpecToProto(spec *keygen.LeveledKeyGenSpec) *KeyGenSpec {
	// Compute the number of keys in each subkey set and total number of subkeys
	levelSizes := make([]int32, spec.NumLevels())
	totalSubKeys := int32(0)
	for i, subKeySet := range spec.SubKeySets {
		levelSizes[i] = int32(len(subKeySet))
		totalSubKeys += levelSizes[i]
	}

	// Build the combined list of all subkeys
	allSubKeys := make([]string, totalSubKeys)
	i := 0
	for _, subKeySet := range spec.SubKeySets {
		for _, subkey := range subKeySet {
			allSubKeys[i] = subkey
			i++
		}
	}

	return &KeyGenSpec{
		LevelSizes: levelSizes,
		SubKeys:    allSubKeys,
		Delim:      spec.Delim,
	}
}

func KeyGenSpecFromProto(protoSpec *KeyGenSpec) *keygen.LeveledKeyGenSpec {
	levelSizes := protoSpec.GetLevelSizes()
	numLevels := len(levelSizes)
	subKeySets := make([][]string, numLevels)
	for i, levelSize := range levelSizes {
		subKeySets[i] = make([]string, levelSize)
	}

	allSubKeys := protoSpec.GetSubKeys()
	k := 0
	for i, subKeySet := range subKeySets {
		for j, _ := range subKeySet {
			subKeySets[i][j] = allSubKeys[k]
			k++
		}
	}

	return &keygen.LeveledKeyGenSpec{
		SubKeySets: subKeySets,
		Delim:      protoSpec.GetDelim(),
	}
}

func TrSpecToProto(spec *strtr.Spec) *TrSpec {
	return &TrSpec{
		Inputs:       spec.Inputs,
		Replacements: spec.Replacements,
	}
}

func TrSpecFromProto(protoSpec *TrSpec) *strtr.Spec {
	return &strtr.Spec{
		Inputs:       protoSpec.GetInputs(),
		Replacements: protoSpec.GetReplacements(),
	}
}

func RangeToProto(r *intgen.Range) *Range {
	return &Range{
		Start: int32(r.Start()),
		Limit: int32(r.Limit()),
	}
}

func RangeFromProto(protoRange *Range) *intgen.Range {
	return intgen.NewRange(int(protoRange.Start), int(protoRange.Limit))
}
