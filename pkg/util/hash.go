package util

import (
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/ripemd160"
)

func Hash256(b []byte) []byte {
	hash1 := sha256.Sum256(b)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:]
}

func Hash160(b []byte) []byte {
	sum := sha256.Sum256(b)
	rip := ripemd160.New()
	io.WriteString(rip, string(sum[:]))
	return rip.Sum(nil)
}

// TODO pass by reference
func ReverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}
