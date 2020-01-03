package message

import (
	"bytes"
	"encoding/binary"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

var ZeroHash = [32]byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

type GetBlocks struct {
	Version            uint32
	HashCount          *common.VarInt
	BlockLocatorHashes [][32]byte
	HashStop           [32]byte
}

func NewGetBlocks(version uint32, blockLocatorHashes [][32]byte, hashStop [32]byte) *GetBlocks {
	var reversedHashStop [32]byte
	copy(reversedHashStop[:], util.ReverseBytes(hashStop[:]))

	length := len(blockLocatorHashes)
	hashCount := common.NewVarInt(uint64(length))
	return &GetBlocks{
		Version:            version,
		HashCount:          hashCount,
		BlockLocatorHashes: blockLocatorHashes,
		HashStop:           reversedHashStop,
	}
}

func (g *GetBlocks) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "getblocks")
	return commandName
}

func (g *GetBlocks) Encode() []byte {
	var version [4]byte
	binary.LittleEndian.PutUint32(version[:4], g.Version)
	hashesBytes := [][]byte{}
	for _, hash := range g.BlockLocatorHashes {
		hashesBytes = append(hashesBytes, hash[:])
	}
	return bytes.Join(
		[][]byte{
			version[:],
			g.HashCount.Encode(),
			bytes.Join(hashesBytes, []byte{}),
			g.HashStop[:],
		},
		[]byte{},
	)
}
