package intgen

import (
	"fmt"
	"math/rand"
)

type Range struct {
	start int
	limit int
}

func NewRange(start, limit int) *Range {
	if limit < start {
		panic(fmt.Sprintf("Bad range creation request(%d, %d): limit cannot be less than start", start, limit))
	}
	return &Range{start, limit}
}

func (r *Range) Start() int {
	return r.start
}

func (r *Range) Limit() int {
	return r.limit
}

func (r *Range) Size() int {
	return r.limit - r.start
}

// Rand will generate a random value in this range
// Will panic if called on an empty range
func (r *Range) Rand() int {
	return r.start + rand.Intn(r.Size())
}

// Ints will generate an integer slice with all the values in the range sequentially
// Will generate empty slice if called on an empty range
func (r *Range) Ints() []int {
	a := make([]int, r.Size())
	for i := 0; i < len(a); i++ {
		a[i] = r.start + i
	}
	return a
}

// Partition will attempt to split this range into n ranges
// At most Size() partitions will be created
func (r *Range) Partition(p int) []*Range {
	partitionSizes := Partition(r.Size(), p)

	output := make([]*Range, len(partitionSizes))
	partitionStart := r.start
	for i, partitionSize := range partitionSizes {
		partitionLimit := partitionStart + partitionSize
		output[i] = NewRange(partitionStart, partitionLimit)
		partitionStart = partitionLimit
	}
	return output
}

// Duplicate will create n copies of the is range
func (r *Range) Duplicate(n int) []*Range {
	output := make([]*Range, n)
	for i := 0; i < len(output); i++ {
		output[i] = r
	}
	return output
}

func (r *Range) String() string {
	return fmt.Sprintf("[%d, %d)", r.start, r.limit)
}
