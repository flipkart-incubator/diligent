package keygen

import "fmt"

type LeveledKeyGen struct {
	subKeySets        [][]string
	delim             string
	subKeyLen         int
	numLevels         int
	numKeys           int
	blockSizePerLevel []int
}

func NewLeveledKeyGen(spec *LeveledKeyGenSpec) *LeveledKeyGen {
	if spec == nil {
		panic("Spec must not be nil")
	}

	if !spec.IsValid() {
		panic("Spec validation failed")
	}

	subKeySets := spec.SubKeySets
	numLevels := len(subKeySets)
	numKeys := 1
	subKeyLen := len(subKeySets[0][0])

	// Loop over set of keys in each level
	for _, subKeySet := range subKeySets {
		n := len(subKeySet)
		numKeys *= n
	}

	return &LeveledKeyGen{
		subKeyLen:         subKeyLen,
		numLevels:         numLevels,
		numKeys:           numKeys,
		blockSizePerLevel: calculateBlockSizePerLevel(subKeySets),
		subKeySets:        subKeySets,
		delim:             defaultDelim,
	}
}

// NumLevels returns the number of levels in the LeveledKeyGen
func (kg *LeveledKeyGen) NumLevels() int {
	return kg.numLevels
}

// NumKeys returns the number of keys in the LeveledKeyGen
func (kg *LeveledKeyGen) NumKeys() int {
	return kg.numKeys
}

// KeyLength returns the size of each key generated by the LeveledKeyGen
func (kg *LeveledKeyGen) KeyLength() int {
	return kg.LengthOfKeyPrefixAtLevel(kg.NumLevels() - 1)
}

// LengthOfKeyPrefixAtLevel returns the size of a key prefix with subkeys upto level n
// n varies in the range [0, spec.NumLevels()]
func (kg *LeveledKeyGen) LengthOfKeyPrefixAtLevel(n int) int {
	if n < 0 || n >= kg.NumLevels() {
		panic(fmt.Sprintf("Invalid level: %d", n))
	}

	// (subkeys * number of subkeys) + number of delimeters
	return (kg.subKeyLen * (n + 1)) + n
}

// BlockSizeAtLevel returns number of keys that start with any specific prefix at level n
// That is, if the keys are A_1, A_2, A_3, B_1, B_2, B_3 then
// BlockSizeAtLevel(0) is 3 (given any prefix at level 0 like "A" or "B", we have will get 3 keys starting with it)
// BlockSizeAtLevel(1) is 1 (given any prefix at level 1 like "A_1" or "B_3", we have will get 1 keys starting with it)
func (kg *LeveledKeyGen) BlockSizeAtLevel(n int) int {
	if n < 0 || n >= kg.NumLevels() {
		panic(fmt.Sprintf("Invalid level: %d", n))
	}
	return kg.blockSizePerLevel[n]
}

// SubKeySets returns all the subkey sets that are used to construct the keys of the LeveledKeyGen
func (kg *LeveledKeyGen) SubKeySets() [][]string {
	return kg.subKeySets
}

// SubKeySet returns the subkey set at a specific level of the LeveledKeyGen
func (kg *LeveledKeyGen) SubKeySet(n int) []string {
	if n < 0 || n >= kg.NumLevels() {
		panic(fmt.Sprintf("Invalid level: %d", n))
	}
	return kg.subKeySets[n]
}

// Key returns the nth key of the LeveledKeyGen in the logical generation order
func (kg *LeveledKeyGen) Key(n int) *LeveledKey {
	if n >= kg.NumKeys() {
		panic(fmt.Sprintf("Attempt to access key at index %d in keyset with %d keys", n, kg.NumKeys()))
	}

	subKeyIndices := make([]int, kg.NumLevels())
	for i := 0; i < kg.NumLevels(); i++ {
		subKeyIndices[i] = n / kg.BlockSizeAtLevel(i)
		n %= kg.BlockSizeAtLevel(i)
	}

	return kg.makeKeyFromIndices(subKeyIndices)
}


// For any prefix at a particular level, how many keys have that prefix. Generate this for every level
func calculateBlockSizePerLevel(subKeySets [][]string) []int {
	// Build slice with number of subkeys in each level
	numLevels := len(subKeySets)
	subKeySetSizes := make([]int, numLevels)
	for i := 0; i < numLevels; i++ {
		subKeySetSizes[i] = len(subKeySets[i])
	}

	// Calculate block size at each level
	p := 1
	output := make([]int, numLevels)
	for i := len(output) - 1; i >= 0; i-- {
		output[i] = p
		p = p * subKeySetSizes[i]
	}
	return output
}

// Construct a generated key given the indices of the subkeys to use from each level
func (kg *LeveledKeyGen) makeKeyFromIndices(indices []int) *LeveledKey {
	if len(indices) != kg.NumLevels() {
		panic(fmt.Sprintf("Attempt to construct keys with wrong number of subkey indices: %d. Expected %d",
			len(indices), kg.NumLevels()))
	}

	subkeys := make([]string, kg.NumLevels())
	for i, k := range indices {
		subkeys[i] = kg.SubKeySets()[i][k]
	}
	return NewLeveledKey(subkeys, defaultDelim)
}

// Exhaustively generates all possible keys - for internal testing use only
func (kg *LeveledKeyGen) allKeys() []*LeveledKey {
	subKeyIndices := make([]int, kg.NumLevels())
	collector := make([]*LeveledKey, kg.NumKeys())
	collectorPos := 0

	kg.generateKeys(0, subKeyIndices, collector, &collectorPos)
	return collector
}

// Recursive helper method for generating all possible keys - for internal testing use only
func (kg *LeveledKeyGen) generateKeys(currentLevel int, subKeyIndices []int, collector []*LeveledKey, collectorPos *int) {
	numSubKeysInCurrentLevel := len(kg.SubKeySet(currentLevel))
	isLastLevel := currentLevel == kg.NumLevels()-1

	for i := 0; i < numSubKeysInCurrentLevel; i++ {
		subKeyIndices[currentLevel] = i

		if isLastLevel {
			collector[*collectorPos] = kg.makeKeyFromIndices(subKeyIndices)
			*collectorPos += 1
		} else {
			kg.generateKeys(currentLevel+1, subKeyIndices, collector, collectorPos)
		}
	}
}