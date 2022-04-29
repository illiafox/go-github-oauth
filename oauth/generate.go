package oauth

import (
	"crypto/rand"
	"encoding/hex"
)

const length = 64

// Generate returns random key (needs optimization)
func Generate() string {
	key := make([]byte, length)
	n, err := rand.Read(key)
	for err != nil {
		n, err = rand.Read(key[n+1:])
	}

	return hex.EncodeToString(key)
}
