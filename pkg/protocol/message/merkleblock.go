package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

// BlockHash return hash of this merkleblock.
// hash256 of version to nonce.
func (m *MerkleBlock) GetBlockHash() [32]byte {
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
	return res
}

func DecodeMerkleBlock(b []byte) (*MerkleBlock, error) {
	if len(b) < 84 {
		return nil, fmt.Errorf("Decode merkle block failed, invalid input: %v", b)
	}
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

	hasheCount, err := common.DecodeVarInt(b)
	if err != nil {
		return nil, err
	}
	hashes := [][32]byte{}
	b = b[len(hasheCount.Encode()):]
	for i := 0; uint64(i) < hasheCount.Data; i++ {
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
		HashCount:         hasheCount,
		Hashes:            hashes,
		FlagBytes:         flagBytes,
		Flags:             flags,
	}, nil
}

// FlagBits return flags of the merkleblock.
func (m *MerkleBlock) FlagBits() []bool {
	res := []bool{}
	for _, flagByte := range m.Flags {
		byteInt := uint8(flagByte)
		for i := 0; i < 8; i++ {
			if (byteInt/uint8(math.Exp2(float64(i))))%uint8(2) == 0x01 {
				res = append(res, true)
			} else {
				res = append(res, false)
			}
		}
	}
	return res
}

func (m *MerkleBlock) Validate() [][32]byte {
	hashes := m.Hashes
	flags := m.FlagBits()
	height := int(math.Ceil(math.Log2(float64(m.TotalTransactions))))
	// マークルパスからrootを計算して m.merkleRoot と一致するか計算

	matchedTxs := [][32]byte{}
	rootHash := calcHash(&hashes, &flags, height, 0, int(m.TotalTransactions), &matchedTxs)
	if bytes.Equal(rootHash[:], m.MerkleRoot[:]) {
		return matchedTxs
	}
	return [][32]byte{}
}

// https://bitcoin.org/en/developer-reference#merkleblock
func calcHash(hashes *[][32]byte, flags *[]bool, height int, pos int, totalTransactions int, matchedTxs *[][32]byte) [32]byte {
	if !(*flags)[0] {
		// フラグが0のとき
		// 先頭のハッシュをこのノードのtxId/ハッシュとする、これより下のノードは探索しない
		*flags = (*flags)[1:]
		h := (*hashes)[0]
		*hashes = (*hashes)[1:]
		return h
	}
	if height == 0 {
		// フラグが1で高さ0(葉ノード)の場合
		// 先頭のハッシュをこのノードのtxIdとして、このトランザクションはマッチ
		*flags = (*flags)[1:]
		h := (*hashes)[0]
		*hashes = (*hashes)[1:]
		*matchedTxs = append(*matchedTxs, h)
		return h
	}
	// calculate left hash
	*flags = (*flags)[1:]
	left := calcHash(hashes, flags, height-1, pos*2, totalTransactions, matchedTxs)
	// calculate right hash if not beyond the end of the array - copy left hash otherwise
	var right [32]byte
	if pos*2+1 < calcTreeWidth(uint(height-1), totalTransactions) {
		right = calcHash(hashes, flags, height-1, pos*2+1, totalTransactions, matchedTxs)
	} else {
		copy(right[:], left[:])
	}
	// combine subhashes
	hash := util.Hash256(bytes.Join([][]byte{left[:], right[:]}, []byte{}))
	var res [32]byte
	copy(res[:], hash)
	return res
}

func calcTreeWidth(height uint, totalTransactions int) int {
	// refer: https://github.com/bitcoin/bitcoin/blob/5961b23898ee7c0af2626c46d5d70e80136578d3/src/merkleblock.h#L65-L68
	return (totalTransactions + (1 << height) - 1) >> height
}

type Merkleblocks struct {
	blocks []*MerkleBlock
}

func NewMerkleBlocks() *Merkleblocks {
	return &Merkleblocks{
		[]*MerkleBlock{},
	}
}

func (m *Merkleblocks) Add(b *MerkleBlock) {
	m.blocks = append(m.blocks, b)
}

func (m *Merkleblocks) Size() int {
	return len(m.blocks)
}

func (m *Merkleblocks) LatestBlock() *MerkleBlock {
	var latest *MerkleBlock
	for _, block := range m.blocks {
		if latest == nil || block.Timestamp < latest.Timestamp {
			latest = block
		}
	}
	return latest
}
