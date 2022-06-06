package keygen

import "strings"

// LeveledKey represents a single key generated as a combination of subkeys from different levels
type LeveledKey struct {
	subKeys []string
	delim string
	key string
}

func NewLeveledKey(subKeys []string, delim string) *LeveledKey {
	key := strings.Join(subKeys, delim)
	return &LeveledKey{
		subKeys: subKeys,
		key: key,
		delim: delim,
	}
}

func (k *LeveledKey) NumLevels() int {
	return len(k.subKeys)
}

func (k *LeveledKey) String() string {
	return k.key
}

func (k *LeveledKey) Prefix(level int) string {
	return strings.Join(k.subKeys[0:level+1], k.delim)
}

