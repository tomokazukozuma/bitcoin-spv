package message

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/spaolacci/murmur3"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type Filterload struct {
	Count      *common.VarInt
	Filter     []byte
	NHashFuncs uint32
	NTweak     uint32
	NFlags     uint8
}

func NewFilterload(size uint32, nHashFuncs uint32, queries [][]byte) *Filterload {
	byteArray := make([]byte, size)
	nTweak := make([]byte, 4)
	for i := 0; i < cap(nTweak); i++ {
		nTweak[i] = util.RandInt8(0, math.MaxUint8)
	}
	nTweakUint32 := binary.BigEndian.Uint32(nTweak)
	for _, query := range queries {
		for i := 0; uint32(i) < nHashFuncs; i++ {
			// 0xFBA4C795 comes from here
			// https://github.com/bitcoin/bitcoin/blob/5961b23898ee7c0af2626c46d5d70e80136578d3/src/bloom.cpp#L52-L56
			seed := uint32(i)*0xFBA4C795 + nTweakUint32
			hashValue := murmur3.Sum32WithSeed(query, seed)
			adjustHashValue := hashValue % (size * uint32(8))
			idx := adjustHashValue >> 3
			value := 1 << (uint32(7) & hashValue)
			byteArray[idx] = byte(value)
		}
	}
	return &Filterload{
		Count:      common.NewVarInt(uint64(size)),
		Filter:     byteArray,
		NHashFuncs: nHashFuncs,
		NTweak:     nTweakUint32,
		NFlags:     uint8(1),
	}
}

// CommandName return message's command name.
func (f *Filterload) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "filterload")
	return commandName
}

// Encode encode message to byte slice.
func (f *Filterload) Encode() []byte {
	nHashFuncsByte := make([]byte, 4)
	nTweakByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(nHashFuncsByte, f.NHashFuncs)
	binary.LittleEndian.PutUint32(nTweakByte, f.NTweak)
	return bytes.Join([][]byte{
		f.Count.Encode(),
		f.Filter,
		nHashFuncsByte,
		nTweakByte,
		[]byte{f.NFlags},
	}, []byte{})
}
