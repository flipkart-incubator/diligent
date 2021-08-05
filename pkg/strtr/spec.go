package strtr

import (
	"fmt"
	"log"
	"math/rand"
)

type Spec struct {
	Inputs       string
	Replacements string
}

func NewRandomTrSpec(inputs string) *Spec {
	ra := []rune(inputs)
	replacements := inputs

	// In this loop we guard against the very rare chance that the shuffle may leave us with
	// Inputs == Replacements. Then we repeat the shuffle
	for replacements == inputs {
		rand.Shuffle(len(ra), func(i, j int) {
			ra[i], ra[j] = ra[j], ra[i]
		})
		replacements = string(ra)
	}

	return NewTrSpec(inputs, replacements)
}

func NewTrSpec(inputs string, replacements string) *Spec {
	tr := &Spec{
		Inputs:       inputs,
		Replacements: replacements,
	}

	if !tr.IsValid() {
		log.Panic(fmt.Sprintf("Invalid transalation spec: %s->%s", inputs, replacements))
	}

	return tr
}

func (spec *Spec) IsValid() bool {
	ra1 := []rune(spec.Inputs)
	ra2 := []rune(spec.Replacements)
	if len(ra1) != len(ra2) {
		return false
	}
	return true
}