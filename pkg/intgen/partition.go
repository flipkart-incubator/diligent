package intgen

import (
	"fmt"
	"math"
)

// Partition divides an integer value n into p partitions of roughly equal size
// It returns the size of each partition
// The number of partitions cannot be greater than n. If greater it will get capped.
func Partition(n, p int) []int {
	if n < 0 {
		panic(fmt.Errorf("trying to partition a negative value: %d", n))
	}
	if p < 0 {
		panic(fmt.Errorf("number of partitions cannot be negative: %d", p))
	}

	// Since there are n distinct values, the number of partitions cannot be more than that
	numPartitions := int(math.Min(float64(n), float64(p)))

	// Deal with empty case
	if numPartitions == 0 {
		return []int{}
	}

	output := make([]int, numPartitions)

	// Get the remainder and quotient of dividing n by numPartitions
	q := n / numPartitions
	r := n % numPartitions

	// Now exactly r partitions should be of size q+1, others should be of size q
	for i := 0; i < r; i++ {
		output[i] = q + 1
	}
	for i := r; i < numPartitions; i++ {
		output[i] = q
	}
	return output
}
