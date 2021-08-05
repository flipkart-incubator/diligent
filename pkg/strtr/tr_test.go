package strtr

import (
	"github.com/flipkart-incubator/diligent/pkg/charset"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test a specific replacement
func TestSpecificTr1(t *testing.T) {
	replacements := "QWERTYUIOPLKJHGFDSAZXCVBNM"
	tr := NewTr(NewTrSpec(charset.AlphaUp, replacements))
	out := tr.Apply(charset.AlphaUp)
	assert.Equal(t, replacements, out)
}

// Test a specific replacement with ignored characters
func TestSpecificTr2(t *testing.T) {
	tr := NewTr(NewTrSpec("ABCD", "WXYZ"))
	out := tr.Apply("CD_AB!!?")
	assert.Equal(t, "YZ_WX!!?", out)
}

// Test a random replacement
func TestRandomTr(t *testing.T) {
	trs := NewRandomTrSpec(charset.AlphaUp)
	assert.NotEqual(t, trs.Inputs, trs.Replacements)

	tr := NewTr(trs)
	out := tr.Apply(charset.AlphaUp)
	assert.Equal(t, trs.Replacements, out)
}
