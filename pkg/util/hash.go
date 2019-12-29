package util

import (
	"crypto/sha256"
	"io"
	"math/rand"
	"time"

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

func RandInt8(min int, max int) uint8 {
	rand.Seed(time.Now().UTC().UnixNano())
	return uint8(min + rand.Intn(max-min))
}
