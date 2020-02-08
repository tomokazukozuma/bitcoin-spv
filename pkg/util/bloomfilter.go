package util

import (
	"github.com/spaolacci/murmur3"
)

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
