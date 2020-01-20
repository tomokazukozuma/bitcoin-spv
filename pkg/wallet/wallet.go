package wallet

import (
	"bytes"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/script"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type Wallet struct {
	Key   *util.Key
	Utxos []*message.Utxo
}

func NewWallet() *Wallet {
	return &Wallet{
		Key:   util.NewKey(),
		Utxos: []*message.Utxo{},
	}
}

func (w *Wallet) GetPublicKey() []byte {
	return w.Key.PublicKey.SerializeUncompressed()
}

func (w *Wallet) GetPublicKeyHash() []byte {
	return util.Hash160(w.GetPublicKey())
}

func (w *Wallet) GetAddress() string {
	return util.EncodeAddress(w.GetPublicKey())
}

func (w *Wallet) AddUtxo(utxo *message.Utxo) {
	alreadyExists := false
	for _, u := range w.Utxos {
		if u.Hash == utxo.Hash && u.N == utxo.N {
			alreadyExists = true
		}
	}
	if alreadyExists {
		return
	}
	w.Utxos = append(w.Utxos, utxo)
}

func (w *Wallet) GetBalance() uint64 {
	var balance uint64
	for _, v := range w.Utxos {
		balance += v.TxOut.Value
	}
	return balance
}

func (w *Wallet) CreateTx(toAddress string, value uint64) *message.Tx {
	// TODO valueに必要なutxoを取得するように修正
	utxos := w.Utxos
	fee := util.CalculateFee(10, len(utxos))
	chargeValue := w.GetBalance() - value - fee
	txouts := w.CreateTxOuts(toAddress, value, chargeValue)
	txins, err := w.CreateTxIns(utxos, txouts)
	if err != nil {
		log.Fatalf("createTxIn: %+v", err)
	}
	// TODO 使用したutxoを削除

	return message.NewTx(uint32(1), txins, txouts, uint32(0)).(*message.Tx)
}

func (w *Wallet) CreateTxOuts(toAddress string, value, chargeValue uint64) []*message.TxOut {
	var txout []*message.TxOut
	lockingScript1 := script.CreateLockingScriptForPKH(util.DecodeAddress(toAddress))
	txout = append(txout, &message.TxOut{
		Value:         value,
		LockingScript: common.NewVarStr(lockingScript1),
	})

	lockingScript2 := script.CreateLockingScriptForPKH(util.DecodeAddress(w.GetAddress()))
	txout = append(txout, &message.TxOut{
		Value:         chargeValue,
		LockingScript: common.NewVarStr(lockingScript2),
	})
	return txout
}

func (w *Wallet) CreateTxIns(utxos []*message.Utxo, txouts []*message.TxOut) ([]*message.TxIn, error) {
	var txins []*message.TxIn
	for _, utxo := range utxos {
		txin := &message.TxIn{
			PreviousOutput: &message.OutPoint{
				Hash: utxo.Hash,
				N:    utxo.N,
			},
			UnlockingScript: utxo.TxOut.LockingScript,
			Sequence:        0xFFFFFFFF,
		}

		tx := message.NewTx(
			uint32(1),
			[]*message.TxIn{txin},
			txouts,
			uint32(0),
		)

		sigHash := util.Hash256(bytes.Join([][]byte{
			tx.Encode(),
			[]byte{0x01, 0x00, 0x00, 0x00},
		}, []byte{}))

		signature, err := w.Sign(sigHash)
		if err != nil {
			return nil, err
		}
		log.Printf("signature len: %+v", len(signature))
		hashType := []byte{0x01}
		signatureWithType := bytes.Join([][]byte{signature, hashType}, []byte{})
		txin.UnlockingScript = script.CreateUnlockingScriptForPKH(signatureWithType, w.GetPublicKey())
		txins = append(txins, txin)
	}
	log.Printf("==== txins len: %+v", len(txins))
	return txins, nil
}

func (w *Wallet) Sign(sigHash []byte) ([]byte, error) {
	return w.Key.Sign(sigHash)
}
