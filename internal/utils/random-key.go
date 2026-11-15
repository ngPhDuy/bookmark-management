package utils

import (
	"crypto/rand"
	"math/big"
)


// Allowed characters, alphabet + numbers
const ALLOWED_CHAR = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate a random key, with random length. Characters are selected from ALLOWED_CHAR
// Use crypto/rand instead of math/rand to avoid seed issues
func RandomKey(length int) (string, error) {
	b := make([]byte, length)
	idx := 0

	for idx < length {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(len(ALLOWED_CHAR))))

		if err != nil {
			return "", err
		}

		b[idx] = ALLOWED_CHAR[i.Int64()]
		idx++
	}

	return string(b), nil
}
