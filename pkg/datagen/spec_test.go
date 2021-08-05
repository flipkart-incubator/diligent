package datagen

import (
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/charset"
	"github.com/flipkart-incubator/diligent/pkg/strgen"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewSpec(t *testing.T) {
	inputRecCounts :=    []int{5, 55, 105, 1005, 10005, 1000005}
	expectedRecCounts := []int{5, 50, 100, 1000, 10000, 1000000}

	for i, recCount := range inputRecCounts {
		spec := NewSpec(recCount, 1024)
		assert.Equal(t, spec.SpecType, specTypeString)
		assert.Equal(t, spec.Version, currentSpecVersion)
		assert.Equal(t, spec.RecordSize, 1024)
		assert.NotNil(t, spec.KeyGenSpec)
		assert.True(t, spec.KeyGenSpec.IsValid())
		assert.Equal(t, spec.KeyGenSpec.NumLevels(), 3)
		assert.Equal(t, spec.KeyGenSpec.NumKeys(), expectedRecCounts[i])
		assert.True(t, spec.UniqTrSpec.IsValid())
		assert.True(t, spec.SmallGrpTrSpec.IsValid())
		assert.True(t, spec.LargeGrpTrSpec.IsValid())
		assert.Equal(t, spec.KeyGenSpec.KeyLength(), len(spec.FixedValue))
	}
}

func TestValidation(t *testing.T) {
	{
		// Bad spec type
		spec := NewSpec(1000, 1024)
		spec.SpecType = "boo"
		assert.False(t, spec.IsValid())
	}
	{
		// Bad version
		spec := NewSpec(1000, 1024)
		spec.Version = 1000
		assert.False(t, spec.IsValid())
	}
	{
		// Bad record size
		spec := NewSpec(1000, 1024)
		spec.RecordSize = -999
		assert.False(t, spec.IsValid())
	}
	{
		// Bad keyset spec
		spec := NewSpec(1000, 1024)
		spec.KeyGenSpec.SubKeySets = [][]string{}
		assert.False(t, spec.IsValid())
	}
	{
		// Bad UniqTrSpec spec
		spec := NewSpec(1000, 1024)
		spec.UniqTrSpec.Replacements = "AA"
		assert.False(t, spec.IsValid())
	}
	{
		// Bad SmallGrpTrSpec spec
		spec := NewSpec(1000, 1024)
		spec.SmallGrpTrSpec.Replacements = "AA"
		assert.False(t, spec.IsValid())
	}
	{
		// Bad LargeGrpTrSpec spec
		spec := NewSpec(1000, 1024)
		spec.LargeGrpTrSpec.Replacements = "AA"
		assert.False(t, spec.IsValid())
	}
	{
		// Bad FixeValue
		spec := NewSpec(1000, 1024)
		spec.FixedValue = ""
		assert.False(t, spec.IsValid())
	}
}

func TestSaveAndRestore(t *testing.T) {
	sg := strgen.NewStrGen(charset.AlphaUp)
	spec1 := NewSpec(1000, 1024)
	tempFileName := fmt.Sprintf("%s/%s", os.TempDir(), sg.RandomString(5))
	fmt.Println(tempFileName)
	err := spec1.SaveToFile(tempFileName)
	assert.Nil(t, err)
	assert.FileExists(t, tempFileName)
	spec2, err := LoadSpecFromFile(tempFileName)
	assert.Nil(t, err)
	assert.Equal(t, spec1, spec2)
	err = os.Remove(tempFileName)
	assert.Nil(t, err)
}
