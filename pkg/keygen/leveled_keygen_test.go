package keygen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicParams1Level(t *testing.T) {
	subKeySets := [][]string{{"A"}}
	spec := NewLeveledKeyGenSpec(subKeySets)
	kg := NewLeveledKeyGen(spec)
	assert.Equal(t, 1, kg.NumLevels())
	assert.Equal(t, 1, kg.NumKeys())
	assert.Equal(t, 1, kg.KeyLength())
	assert.Equal(t, 1, len(kg.Key(0).String()))
	assert.Equal(t, 1, kg.LengthOfKeyPrefixAtLevel(0))
	assert.Equal(t, 1, kg.BlockSizeAtLevel(0))
	assert.Equal(t, subKeySets, kg.SubKeySets())
	assert.Equal(t, subKeySets[0], kg.SubKeySet(0))
}

func TestBasicParams2Level(t *testing.T) {
	subKeySets := [][]string{{"A", "B"}, {"1", "2", "3"}}
	spec := NewLeveledKeyGenSpec(subKeySets)
	kg := NewLeveledKeyGen(spec)
	assert.Equal(t, 2, kg.NumLevels())
	assert.Equal(t, 6, kg.NumKeys())
	assert.Equal(t, 3, kg.KeyLength())
	assert.Equal(t, 3, len(kg.Key(0).String()))
	assert.Equal(t, 1, kg.LengthOfKeyPrefixAtLevel(0))
	assert.Equal(t, 3, kg.LengthOfKeyPrefixAtLevel(1))
	assert.Equal(t, 3, kg.BlockSizeAtLevel(0))
	assert.Equal(t, 1, kg.BlockSizeAtLevel(1))
	assert.Equal(t, subKeySets, kg.SubKeySets())
	assert.Equal(t, subKeySets[0], kg.SubKeySet(0))
	assert.Equal(t, subKeySets[1], kg.SubKeySet(1))
}

func TestBasicParams3Level(t *testing.T) {
	subKeySets := [][]string{{"A", "B"}, {"1", "2", "3"}, {"W", "X", "Y", "Z"}}
	spec := NewLeveledKeyGenSpec(subKeySets)
	kg := NewLeveledKeyGen(spec)
	assert.Equal(t, 3, kg.NumLevels())
	assert.Equal(t, 24, kg.NumKeys())
	assert.Equal(t, 5, kg.KeyLength())
	assert.Equal(t, 5, len(kg.Key(0).String()))
	assert.Equal(t, 1, kg.LengthOfKeyPrefixAtLevel(0))
	assert.Equal(t, 3, kg.LengthOfKeyPrefixAtLevel(1))
	assert.Equal(t, 5, kg.LengthOfKeyPrefixAtLevel(2))
	assert.Equal(t, 12, kg.BlockSizeAtLevel(0))
	assert.Equal(t, 4, kg.BlockSizeAtLevel(1))
	assert.Equal(t, 1, kg.BlockSizeAtLevel(2))
	assert.Equal(t, subKeySets, kg.SubKeySets())
	assert.Equal(t, subKeySets[0], kg.SubKeySet(0))
	assert.Equal(t, subKeySets[1], kg.SubKeySet(1))
	assert.Equal(t, subKeySets[2], kg.SubKeySet(2))
}

// First validate the allKeys() internal method is generating the exhaustive list of keys
// We will use this below to validate the public Keys() method
func TestAllKeys(t *testing.T) {
	spec1 := NewLeveledKeyGenSpec([][]string{{"A"}})
	lkg1 := NewLeveledKeyGen(spec1)
	expectedKeys1 := []*LeveledKey{
		NewLeveledKey([]string{"A"}, "_"),
	}
	assert.Equal(t, expectedKeys1, lkg1.allKeys())

	spec2 := NewLeveledKeyGenSpec([][]string{{"A", "B"}})
	lkg2 := NewLeveledKeyGen(spec2)
	expectedKeys2 := []*LeveledKey{
		NewLeveledKey([]string{"A"}, "_"),
		NewLeveledKey([]string{"B"}, "_"),
	}
	assert.Equal(t, expectedKeys2, lkg2.allKeys())

	spec3 := NewLeveledKeyGenSpec([][]string{{"A", "B"}, {"1", "2", "3"}})
	lkg3 := NewLeveledKeyGen(spec3)
	expectedKeys3 := []*LeveledKey{
		NewLeveledKey([]string{"A", "1"}, "_"),
		NewLeveledKey([]string{"A", "2"}, "_"),
		NewLeveledKey([]string{"A", "3"}, "_"),
		NewLeveledKey([]string{"B", "1"}, "_"),
		NewLeveledKey([]string{"B", "2"}, "_"),
		NewLeveledKey([]string{"B", "3"}, "_"),
	}
	assert.Equal(t, expectedKeys3, lkg3.allKeys())
}

func TestKeyGen(t *testing.T) {
	spec1 := NewLeveledKeyGenSpec([][]string{{"A", "B"}, {"1", "2", "3"}, {"W", "X", "Y", "Z"}})
	lkg1 := NewLeveledKeyGen(spec1)
	expectedKeys1 := lkg1.allKeys()
	assert.Equal(t, 24, len(expectedKeys1))
	for i, k := range expectedKeys1 {
		assert.Equal(t, k, lkg1.Key(i))
	}
}

func TestPanicOnBadStrKeySetCreate(t *testing.T) {
	assert.Panics(t, func() { NewLeveledKeyGen(nil) })
}

func TestPanicOnBadKeyIndex(t *testing.T) {
	assert.Panics(t, func() {
		spec1 := NewLeveledKeyGenSpec([][]string{{"A"}})
		lkg1 := NewLeveledKeyGen(spec1)
		lkg1.Key(-1)
	})
	assert.NotPanics(t, func() {
		spec1 := NewLeveledKeyGenSpec([][]string{{"A"}})
		lkg1 := NewLeveledKeyGen(spec1)
		lkg1.Key(0)
	})
	assert.Panics(t, func() {
		spec1 := NewLeveledKeyGenSpec([][]string{{"A"}})
		lkg1 := NewLeveledKeyGen(spec1)
		lkg1.Key(1)
	})
}
