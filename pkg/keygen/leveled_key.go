package keygen

import "strings"

type LeveledKey struct {
	subKeys []string
	key string
	delim string
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

