package message

import (
	"bytes"
	"encoding/binary"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

type TxOut struct {
	Value         uint64
	LockingScript *common.VarStr
}

func (out *TxOut) Encode() []byte {
	valueBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(valueBytes, out.Value)
	return bytes.Join([][]byte{
		valueBytes,
		out.LockingScript.Encode(),
	}, []byte{})
}

func DecodeTxOut(b []byte) (*TxOut, error) {
	value := binary.LittleEndian.Uint64(b[0:8])
	pkScript, _ := common.DecodeVarStr(b[8:])
	return &TxOut{
		Value:         value,
		LockingScript: pkScript,
	}, nil
}
