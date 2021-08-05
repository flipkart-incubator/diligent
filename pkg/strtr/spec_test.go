package strtr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpecValidation(t *testing.T) {
	spec1 := &Spec{
		Inputs:       "ABC",
		Replacements: "XYZ",
	}
	assert.True(t, spec1.IsValid())

	spec2 := &Spec{
		Inputs:       "ABC",
		Replacements: "XY",
	}
	assert.False(t, spec2.IsValid())
}
