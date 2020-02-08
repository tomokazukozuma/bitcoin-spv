package util

import (
	"encoding/binary"
	"math"
	"math/rand"
	"time"

	"github.com/spaolacci/murmur3"
)

func GenerateNTweak() uint32 {
	nTweak := make([]byte, 4)
	for i := 0; i < cap(nTweak); i++ {
		nTweak[i] = randInt8(0, math.MaxUint8) // 0 - 255
	}
	return binary.BigEndian.Uint32(nTweak)
}

func CreateBloomFilter(byteSize uint32, nHashFuncs uint32, queries [][]byte, nTweak uint32) []byte {
	byteArray := make([]byte, byteSize)
	for _, query := range queries {
		for i := 0; uint32(i) < nHashFuncs; i++ {
			seed := generateSeed(i, nTweak)
			hashValue := murmur3.Sum32WithSeed(query, seed)
			adjustHashValue := hashValue % (byteSize * uint32(8))
			idx := adjustHashValue >> 3
			value := 1 << (uint32(7) & hashValue)
			byteArray[idx] = byte(value)
		}
	}
	return byteArray
}

func generateSeed(i int, nTweak uint32) uint32 {
	return (uint32(i)*0xFBA4C795 + nTweak) & 0xffffffff
}
func randInt8(min int, max int) uint8 {
	rand.Seed(time.Now().UTC().UnixNano())
	return uint8(min + rand.Intn(max-min))
}
