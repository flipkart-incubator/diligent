package intgen

import (
	"math/rand"
)

// Shuffle carries out an in-place shuffle of a given integer slice
func Shuffle(a []int) {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}
