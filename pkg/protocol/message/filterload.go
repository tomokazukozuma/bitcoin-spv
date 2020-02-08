package message

import (
	"bytes"
	"encoding/binary"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type Filterload struct {
	Length     *common.VarInt
	Filter     []byte
	NHashFuncs uint32
	NTweak     uint32
	NFlags     uint8
}

func NewFilterload(size uint32, nHashFuncs uint32, queries [][]byte) protocol.Message {
	nTweak := util.GenerateNTweak()
	return &Filterload{
		Length:     common.NewVarInt(uint64(size)),
		Filter:     util.CreateBloomFilter(size, nHashFuncs, queries, nTweak),
		NHashFuncs: nHashFuncs,
		NTweak:     nTweak,
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
		f.Length.Encode(),
		f.Filter,
		nHashFuncsByte,
		nTweakByte,
		[]byte{f.NFlags},
	}, []byte{})
}
