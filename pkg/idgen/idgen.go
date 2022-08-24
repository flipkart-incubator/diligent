package idgen

import (
	"crypto/sha1"
	"fmt"
	"time"
)

func GenerateId16() string {
	// Generate a sha1 hash of current timestamp, and make a 16 char string from it
	now := time.Now()
	hasher := sha1.New()
	hasher.Write([]byte(now.String()))
	id := fmt.Sprintf("%X", hasher.Sum(nil)[:8])
	return id
}
