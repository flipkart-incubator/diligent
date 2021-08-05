package intgen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartition1(t *testing.T) {
	// Test all cases with a small partition
	assert.Equal(t, []int{}, Partition(6, 0))
	assert.Equal(t, []int{6}, Partition(6, 1))
	assert.Equal(t, []int{3, 3}, Partition(6, 2))
	assert.Equal(t, []int{2, 2, 2}, Partition(6, 3))
	assert.Equal(t, []int{1, 1, 1, 3}, Partition(6, 4))
	assert.Equal(t, []int{1, 1, 1, 1, 2}, Partition(6, 5))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1}, Partition(6, 6))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1}, Partition(6, 7))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1}, Partition(6, 8))

	// Test cases for a large partition
	assert.Equal(t, []int{105, 105, 107}, Partition(317, 3))
}
