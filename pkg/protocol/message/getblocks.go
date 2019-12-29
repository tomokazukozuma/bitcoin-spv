package message

import (
	"bytes"
	"encoding/binary"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

type GetBlocks struct {
	Version            uint32
	HashCount          *common.VarInt
	BlockLocatorHashes [][32]byte
	HashStop           [32]byte
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
