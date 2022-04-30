package message

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/script"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type Tx struct {
	Version    uint32
	TxInCount  *common.VarInt
	TxIn       []*TxIn
	TxOutCount *common.VarInt
	TxOut      []*TxOut
	LockTime   uint32
}

type Utxo struct {
	Hash  [32]byte
	N     uint32
	TxOut *TxOut
}

func NewTx(version uint32, txin []*TxIn, txout []*TxOut, locktime uint32) protocol.Message {
	return &Tx{
		Version:    version,
		TxInCount:  common.NewVarInt(uint64(len(txin))),
		TxIn:       txin,
		TxOutCount: common.NewVarInt(uint64(len(txout))),
		TxOut:      txout,
		LockTime:   locktime,
	}
}

func (tx *Tx) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "tx")
	return commandName
}

func (tx *Tx) ID() [32]byte {
	var res [32]byte
	hash := util.Hash256(tx.Encode())
	for i, j := 0, len(hash)-1; i < j; i, j = i+1, j-1 {
		hash[i], hash[j] = hash[j], hash[i]
	}
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

func (tx *Tx) GetUtxo(pubkeyHash []byte) []*Utxo {
	var utxo []*Utxo
	for index, txout := range tx.TxOut {
		// TODO locking scriptの種類をみて、データ部だけのチェックにする
		if bytes.HasPrefix(txout.LockingScript.Data, script.CreateLockingScriptForPKH(pubkeyHash)) {
			utxo = append(utxo, &Utxo{
				Hash:  tx.ID(),
				N:     uint32(index),
				TxOut: txout,
			})
		}
	}
	return utxo
}
