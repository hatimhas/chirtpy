package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// rand.Read Fills the randomBytes slice with cryptographically secure random data
// rand.Read() returns:
//	n (number of bytes read)
//	err (if any)

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32) // 256 bits = 32 bytes
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(randomBytes), nil
}
