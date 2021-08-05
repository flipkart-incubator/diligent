package keygen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicSpecParams1Level(t *testing.T) {
	subKeySets := [][]string{{"A"}}
	spec := NewLeveledKeyGenSpec(subKeySets)
	assert.Equal(t, 1, spec.NumLevels())
	assert.Equal(t, 1, spec.NumKeys())
	assert.Equal(t, 1, spec.KeyLength())
	assert.Equal(t, 1, spec.LengthOfKeyPrefixAtLevel(0))
	assert.Equal(t, subKeySets, spec.SubKeySets)
}

func TestBasicSpecParams2Level(t *testing.T) {
	subKeySets := [][]string{{"A", "B"}, {"1", "2", "3"}}
	spec := NewLeveledKeyGenSpec(subKeySets)
	assert.Equal(t, 2, spec.NumLevels())
	assert.Equal(t, 6, spec.NumKeys())
	assert.Equal(t, 3, spec.KeyLength())
	assert.Equal(t, 1, spec.LengthOfKeyPrefixAtLevel(0))
	assert.Equal(t, 3, spec.LengthOfKeyPrefixAtLevel(1))
	assert.Equal(t, subKeySets, spec.SubKeySets)
}

func TestBasicSpecParams3Level(t *testing.T) {
	subKeySets := [][]string{{"A", "B"}, {"1", "2", "3"}, {"W", "X", "Y", "Z"}}
	spec := NewLeveledKeyGenSpec(subKeySets)
	assert.Equal(t, 3, spec.NumLevels())
	assert.Equal(t, 24, spec.NumKeys())
	assert.Equal(t, 5, spec.KeyLength())
	assert.Equal(t, 1, spec.LengthOfKeyPrefixAtLevel(0))
	assert.Equal(t, 3, spec.LengthOfKeyPrefixAtLevel(1))
	assert.Equal(t, 5, spec.LengthOfKeyPrefixAtLevel(2))
	assert.Equal(t, subKeySets, spec.SubKeySets)
}

func TestCreateRandom(t *testing.T) {
	testCases := [][]int{
		{1},
		{1, 1, 1},
		{10, 20, 30},
	}

	for _, numSubKeysByLevel := range testCases {
		numLevels := len(numSubKeysByLevel)

		spec := NewRandomLeveledKeyGenSpec(numSubKeysByLevel, 10)
		// Check if every SubKeySet is of expected size
		for i, n := range numSubKeysByLevel {
			assert.Equal(t, n, len(spec.SubKeySets[i]))
		}

		// Check if every SubKeySet has unique strings
		for i := 0; i < numLevels; i++ {
			assert.True(t, isUnique(spec.SubKeySets[i]))
		}

		// Check if every SubKeySet consists of unique strings of expected length
		for i := 0; i < numLevels; i++ {
			// Check every SubKey in the SubKeySet at level i
			for _, s := range spec.SubKeySets[i] {
				assert.Equal(t, 10, len(s))
			}
		}
	}
}

func TestValidation(t *testing.T) {
	// Zero levels
	spec1 := &LeveledKeyGenSpec{SubKeySets: [][]string{}, Delim: defaultDelim}
	assert.False(t, spec1.IsValid())
	// Empty level
	spec2 := &LeveledKeyGenSpec{SubKeySets: [][]string{{"A"}, {}, {"C"}}, Delim: defaultDelim}
	assert.False(t, spec2.IsValid())
	// Different sub key lengths
	spec3 := &LeveledKeyGenSpec{SubKeySets: [][]string{{"A"}, {"B", "CD"}}, Delim: defaultDelim}
	assert.False(t, spec3.IsValid())
}

func TestPanicOnBadSpecCreate(t *testing.T) {
	// Panic on zero levels
	assert.Panics(t, func() {
		NewRandomLeveledKeyGenSpec([]int{}, 10)
	})
	// Panic on level with zero subkeys
	assert.Panics(t, func() {
		NewRandomLeveledKeyGenSpec([]int{5, 0, 5}, 10)
	})
	// Panic on zero subkey length
	assert.Panics(t, func() {
		NewRandomLeveledKeyGenSpec([]int{5}, 0)
	})
}

func isUnique(strings []string) bool {
	set := make(map[string]int)

	for _, s := range strings {
		_, ok := set[s]
		if ok {
			return false
		}
		set[s] = 0
	}
	return true
}
