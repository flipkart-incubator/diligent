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

	// Check if n is exactly divisible by numPartitions
	if (n % numPartitions) == 0 {
		// Yes, so its just numPartitions partitions all of size n/numPartitions
		partitionSize := n / numPartitions
		for i := 0; i < len(output); i++ {
			output[i] = partitionSize
		}
	} else {
		// No, so create numPartitions-1 partitions of size floor(n/numPartitions)
		// then create last partition with remaining size
		partitionSize := n / numPartitions
		for i := 0; i < len(output) - 1; i++ {
			output[i] = partitionSize
		}
		output[len(output) - 1] = n - (partitionSize * (numPartitions-1))
	}

	return output
}
