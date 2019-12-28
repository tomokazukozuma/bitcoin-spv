package util

import "crypto/sha256"

func Hash256(b []byte) []byte {
	hash1 := sha256.Sum256(b)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:]
}
