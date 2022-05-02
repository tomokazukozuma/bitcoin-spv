package message

import (
	"bytes"
	"encoding/binary"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type TxIn struct {
	PreviousOutput  *OutPoint
	UnlockingScript *common.VarStr
	Sequence        uint32
}

type OutPoint struct {
	Hash [32]byte
	N    uint32
}

func (in *TxIn) Encode() []byte {
	sequenceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequenceBytes, in.Sequence)
	return bytes.Join([][]byte{
		in.PreviousOutput.Encode(),
		in.UnlockingScript.Encode(),
		sequenceBytes,
	}, []byte{})
}

func (p *OutPoint) Encode() []byte {
	hash := make([]byte, 32)
	copy(hash, p.Hash[:])
	util.ReverseBytes(hash)

	n := make([]byte, 4)
	binary.LittleEndian.PutUint32(n, p.N)
	return bytes.Join([][]byte{
		hash,
		n,
	}, []byte{})
}

func DecodeTxIn(b []byte) (*TxIn, error) {
	var hash [32]byte
	copy(hash[:], b[0:32])
	util.ReverseBytes(hash[:])
	n := binary.LittleEndian.Uint32(b[32:36])
	out := &OutPoint{
		Hash: hash,
		N:    n,
	}
	b = b[36:]
	signatureScript, err := common.DecodeVarStr(b)
	if err != nil {
		return nil, err
	}
	length := len(signatureScript.Encode())
	b = b[length:]
	sequence := binary.LittleEndian.Uint32(b[:4])
	return &TxIn{
		PreviousOutput:  out,
		UnlockingScript: signatureScript,
		Sequence:        sequence,
	}, nil
}
