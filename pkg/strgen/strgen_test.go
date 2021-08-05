package strgen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	charset := "ABCD"
	sg := NewStrGen(charset)

	// Generate strings of different lengths
	// Check that length is as expected and strings use the specified charset
	for length := 0; length < 10; length++ {
		s := sg.RandomString(length)
		assert.Equal(t, length, len(s))
		assert.True(t, stringHasCharset(s, charset))
	}
}

func TestRandomStringsUnique(t *testing.T) {
	charset := "ABCD"
	sg := NewStrGen(charset)

	for length := 5; length < 10; length++ {
		for num := 0; num < 100; num++ {
			strings, err := sg.RandomStringsUnique(length, num)

			assert.Nil(t, err)
			assert.Equal(t, num, len(strings))
			assert.True(t, isUnique(strings))
			for _, s := range strings {
				assert.True(t, stringHasCharset(s, charset))
			}
		}
	}
}


func stringHasCharset(s string, charset string) bool {
	cmap := make(map[rune]bool, len(charset))
	for _, c := range []rune(charset) {
		cmap[c] = true
	}

	for _, c := range []rune(s) {
		_, ok := cmap[c]
		if !ok {
			return false
		}
	}
	return true
}

func isUnique(strings []string) bool {
	set := make(map[string]bool)

	for _, s := range strings {
		_, ok := set[s]
		if ok {
			return false
		}
		set[s] = true
	}
	return true
}