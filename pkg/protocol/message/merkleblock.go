package message

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type MerkleBlock struct {
	Version           uint32
	PrevBlock         [32]byte
	MerkleRoot        [32]byte
	Timestamp         uint32
	Bits              uint32
	Nonce             uint32
	TotalTransactions uint32
	HashCount         *common.VarInt
	Hashes            [][32]byte
	FlagBytes         *common.VarInt
	Flags             []byte
}

func (g *MerkleBlock) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "merkleblock")
	return commandName
}

func (m *MerkleBlock) GetBlockHash() string {
	var res [32]byte
	versionByte := make([]byte, 4)
	timestampByte := make([]byte, 4)
	bitsByte := make([]byte, 4)
	nonceByte := make([]byte, 4)

	binary.LittleEndian.PutUint32(versionByte, m.Version)
	binary.LittleEndian.PutUint32(timestampByte, m.Timestamp)
	binary.LittleEndian.PutUint32(bitsByte, m.Bits)
	binary.LittleEndian.PutUint32(nonceByte, m.Nonce)

	bs := bytes.Join([][]byte{
		versionByte,
		m.PrevBlock[:],
		m.MerkleRoot[:],
		timestampByte,
		bitsByte,
		nonceByte,
	}, []byte{})

	copy(res[:], util.Hash256(bs))
	util.ReverseBytes(res[:])
	return hex.EncodeToString(res[:])
}

func DecodeMerkleBlock(b []byte) (*MerkleBlock, error) {
	version := binary.LittleEndian.Uint32(b[0:4])
	var prevBlockArr [32]byte
	var merkleRootArr [32]byte
	copy(prevBlockArr[:], b[4:36])
	copy(merkleRootArr[:], b[36:68])
	timestamp := binary.LittleEndian.Uint32(b[68:72])
	bits := binary.LittleEndian.Uint32(b[72:76])
	nonce := binary.LittleEndian.Uint32(b[76:80])
	totalTransactions := binary.LittleEndian.Uint32(b[80:84])

	b = b[84:]

	hashCount, err := common.DecodeVarInt(b)
	if err != nil {
		return nil, err
	}
	var hashes [][32]byte
	b = b[len(hashCount.Encode()):]
	for i := 0; uint64(i) < hashCount.Data; i++ {
		var byteArray [32]byte
		copy(byteArray[:], b[:32])
		b = b[32:]
		hashes = append(hashes, byteArray)
	}

	flagBytes, err := common.DecodeVarInt(b)
	if err != nil {
		return nil, err
	}
	b = b[len(flagBytes.Encode()):]
	flags := b[:flagBytes.Data]

	return &MerkleBlock{
		Version:           version,
		PrevBlock:         prevBlockArr,
		MerkleRoot:        merkleRootArr,
		Timestamp:         timestamp,
		Bits:              bits,
		Nonce:             nonce,
		TotalTransactions: totalTransactions,
		HashCount:         hashCount,
		Hashes:            hashes,
		FlagBytes:         flagBytes,
		Flags:             flags,
	}, nil
}

func (m *MerkleBlock) Validate() [][32]byte {
	hashes := m.Hashes
	flags := m.decodeFlagBits()
	height := int(math.Ceil(math.Log2(float64(m.TotalTransactions))))

	var matchedTxs [][32]byte
	rootHash := calcHash(flags, height, 0, int(m.TotalTransactions), hashes, matchedTxs)
	if bytes.Equal(rootHash[:], m.MerkleRoot[:]) {
		return matchedTxs
	}
	return [][32]byte{}
}

func (m *MerkleBlock) decodeFlagBits() (flags []bool) {
	for _, flagByte := range m.Flags {
		byteInt := flagByte
		for i := 0; i < 8; i++ {
			if (byteInt/uint8(math.Exp2(float64(i))))%uint8(2) == 0x01 {
				flags = append(flags, true)
			} else {
				flags = append(flags, false)
			}
		}
	}
	return
}

func calcHash(flags []bool, height, pos, totalTransactions int, hashes, matchedTxs [][32]byte) [32]byte {
	if !flags[0] {
		flags = flags[1:]
		h := (hashes)[0]
		hashes = hashes[1:]
		return h
	}
	if height == 0 {
		flags = flags[1:]
		h := hashes[0]
		hashes = hashes[1:]
		matchedTxs = append(matchedTxs, h)
		return h
	}

	flags = flags[1:]
	left := calcHash(flags, height-1, pos*2, totalTransactions, hashes, matchedTxs)
	var right [32]byte
	if pos*2+1 < calcTreeWidth(uint(height-1), totalTransactions) {
		right = calcHash(flags, height-1, pos*2+1, totalTransactions, hashes, matchedTxs)
	} else {
		copy(right[:], left[:])
	}
	hash := util.Hash256(bytes.Join([][]byte{left[:], right[:]}, []byte{}))
	var res [32]byte
	copy(res[:], hash)
	return res
}

func calcTreeWidth(height uint, totalTransactions int) int {
	return (totalTransactions + (1 << height) - 1) >> height
}
