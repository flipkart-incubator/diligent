package intgen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShuffle(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	Shuffle(b)
	assert.ElementsMatch(t, a, b)
	assert.NotEqual(t, a, b) // Note: There is a rare chance of a==b
}