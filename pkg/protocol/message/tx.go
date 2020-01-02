package message

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

type Tx struct {
	Version    uint32
	TxInCount  *common.VarInt
	TxIn       []*TxIn
	TxOutCount *common.VarInt
	TxOut      []*TxOut
	LockTime   uint32
}

type TxIn struct {
	PreviousOutput  *OutPoint
	SignatureScript *common.VarStr
	Sequence        uint32
}

type OutPoint struct {
	Hash  [32]byte
	Index uint32
}

type TxOut struct {
	Value         uint64
	LockingScript *common.VarStr
}

func NewTx() *Tx {
	return &Tx{
		Version:    0,
		TxInCount:  nil,
		TxIn:       nil,
		TxOutCount: nil,
		TxOut:      nil,
		LockTime:   0,
	}
}

func (tx *Tx) ID() [32]byte {
	var res [32]byte
	hash := util.Hash256(tx.Encode())
	copy(res[:], hash)
	return res
}

func (tx *Tx) Encode() []byte {
	versionBytes := make([]byte, 4)
	lockTimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(versionBytes, tx.Version)
	binary.LittleEndian.PutUint32(lockTimeBytes, tx.LockTime)

	txInBytes := [][]byte{}
	for _, in := range tx.TxIn {
		txInBytes = append(txInBytes, in.Encode())
	}

	txOutBytes := [][]byte{}
	for _, out := range tx.TxOut {
		txOutBytes = append(txOutBytes, out.Encode())
	}

	return bytes.Join([][]byte{
		versionBytes,
		tx.TxInCount.Encode(),
		bytes.Join(txInBytes, []byte{}),
		tx.TxOutCount.Encode(),
		bytes.Join(txOutBytes, []byte{}),
		lockTimeBytes,
	}, []byte{})
}

func DecodeTx(b []byte) (*Tx, error) {
	version := binary.LittleEndian.Uint32(b[0:4])
	b = b[4:]

	var txIns []*TxIn
	txInCount, err := common.DecodeVarInt(b)
	if err != nil {
		return nil, err
	}
	b = b[len(txInCount.Encode()):]
	for i := 0; uint64(i) < txInCount.Data; i++ {
		txIn, err := DecodeTxIn(b)
		if err != nil {
			return nil, err
		}
		txIns = append(txIns, txIn)
		len := len(txIn.Encode())
		b = b[len:]
	}

	var txOuts []*TxOut
	txOutCount, err := common.DecodeVarInt(b)
	if err != nil {
		return nil, err
	}
	b = b[len(txOutCount.Encode()):]
	for i := 0; uint64(i) < txOutCount.Data; i++ {
		txOut, err := DecodeTxOut(b)
		if err != nil {
			return nil, err
		}
		txOuts = append(txOuts, txOut)
		len := len(txOut.Encode())
		b = b[len:]
	}
	if len(b) != 4 {
		return nil, fmt.Errorf("decode Transaction failed, invalid input: %v", b)
	}
	lockTime := binary.LittleEndian.Uint32(b[0:4])
	return &Tx{
		Version:    version,
		TxInCount:  txInCount,
		TxIn:       txIns,
		TxOutCount: txOutCount,
		TxOut:      txOuts,
		LockTime:   lockTime,
	}, nil
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

func DecodeTxOut(b []byte) (*TxOut, error) {
	value := binary.LittleEndian.Uint64(b[0:8])
	pkScript, _ := common.DecodeVarStr(b[8:])
	return &TxOut{
		Value:         value,
		LockingScript: pkScript,
	}, nil
}

func (p *OutPoint) Encode() []byte {
	indexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexBytes, p.Index)
	return bytes.Join([][]byte{
		p.Hash[:],
		indexBytes,
	}, []byte{})

}

func (out *TxOut) Encode() []byte {
	valueBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(valueBytes, out.Value)
	return bytes.Join([][]byte{
		valueBytes,
		out.LockingScript.Encode(),
	}, []byte{})
}
