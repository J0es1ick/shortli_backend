package shortener

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateShortCode(originalURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(originalURL))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	shortCode := hashString[:12]

	return shortCode
}
