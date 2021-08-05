package strgen

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type StrGen struct {
	charset []rune
}

func NewStrGen(charset string) *StrGen {
	rand.Seed(time.Now().UnixNano())
	return &StrGen{
		charset: []rune(charset),
	}
}

// RandomString will generate a random string of a given fixed length
func (g *StrGen) RandomString(length int) string {
	if length < 0 {
		panic("length of strings cannot be negative")
	}

	ra := make([]rune, length)
	for i := 0; i < length; i++ {
		p := rand.Intn(len(g.charset))
		ra[i] = g.charset[p]
	}

	return string(ra)
}

// RandomStringsUnique will generate the required number of random strings, each of a given fixed length
func (g *StrGen) RandomStringsUnique(length int, num int) ([]string, error) {
	if length < 0 {
		panic("length of strings cannot be negative")
	}
	if num < 0 {
		panic("number of strings cannot be negative")
	}

	numSymbols := float64(len(g.charset))
	universe := math.Pow(numSymbols, float64(length))
	if float64(num) > universe {
		return nil, fmt.Errorf("not possible to generate %d unique fixed length strings of length %d", num, length)
	}

	// Generate random strings and put them in a map till we have the required number of keys
	set := make(map[string]bool)
	for len(set) < num {
		set[g.RandomString(length)] = true
	}

	// Take the keys out of the map, put them in a slice and return
	output := make([]string, 0, num)
	for s := range set {
		output = append(output, s)
	}

	return output, nil
}



