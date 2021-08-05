package keygen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLeveledStrKeys(t *testing.T) {
	k1 := NewLeveledKey([]string{"A"}, "_")
	assert.Equal(t, 1, k1.NumLevels())
	assert.Equal(t, "A", k1.String())
	assert.Equal(t, "A", k1.Prefix(0))

	k2 := NewLeveledKey([]string{"A", "B"}, "_")
	assert.Equal(t, 2, k2.NumLevels())
	assert.Equal(t, "A_B", k2.String())
	assert.Equal(t, "A", k2.Prefix(0))
	assert.Equal(t, "A_B", k2.Prefix(1))

	k3 := NewLeveledKey([]string{"A", "B", "C"}, "_")
	assert.Equal(t, 3, k3.NumLevels())
	assert.Equal(t, "A_B_C", k3.String())
	assert.Equal(t, "A", k3.Prefix(0))
	assert.Equal(t, "A_B", k3.Prefix(1))
	assert.Equal(t, "A_B_C", k3.Prefix(2))
}
