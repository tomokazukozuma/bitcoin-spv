package message

import (
	"bytes"
	"encoding/binary"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

type TxIn struct {
	PreviousOutput  *OutPoint
	SignatureScript *common.VarStr
	Sequence        uint32
}

type OutPoint struct {
	Hash  [32]byte
	Index uint32
}

func (in *TxIn) Encode() []byte {
	sequenceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequenceBytes, in.Sequence)
	return bytes.Join([][]byte{
		in.PreviousOutput.Encode(),
		in.SignatureScript.Encode(),
		sequenceBytes,
	}, []byte{})
}

func (p *OutPoint) Encode() []byte {
	indexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexBytes, p.Index)
	return bytes.Join([][]byte{
		p.Hash[:],
		indexBytes,
	}, []byte{})
}

func DecodeTxIn(b []byte) (*TxIn, error) {
	var hash [32]byte
	copy(hash[:], b[0:32])
	index := binary.LittleEndian.Uint32(b[32:36])
	out := &OutPoint{
		Hash:  hash,
		Index: index,
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
		SignatureScript: signatureScript,
		Sequence:        sequence,
	}, nil
}
