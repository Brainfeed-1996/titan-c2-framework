package transport

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID creates a random ID for agents/jobs
func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
