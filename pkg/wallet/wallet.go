package wallet

import (
	"bytes"
	"log"
	"sort"

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
	// TODO TxInで使われいないかチェック
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
	utxos := w.getEnoughUtxos(value)
	fee := util.CalculateFee(10, len(utxos))
	txouts := w.CreateTxOuts(toAddress, value, fee)
	txins, err := w.CreateTxIns(utxos, txouts)
	if err != nil {
		log.Fatalf("createTxIn: %+v", err)
	}
	for _, utxo := range utxos {
		w.removeUtxo(utxo)
	}
	return message.NewTx(uint32(1), txins, txouts, uint32(0)).(*message.Tx)
}

func (w *Wallet) getEnoughUtxos(value uint64) (utxos []*message.Utxo) {
	sort.Slice(w.Utxos, func(i, j int) bool { return w.Utxos[i].TxOut.Value > w.Utxos[j].TxOut.Value })
	var totalVAlue uint64
	for _, utxo := range w.Utxos {
		utxos = append(utxos, utxo)
		totalVAlue += utxo.TxOut.Value
		if value <= totalVAlue {
			return
		}
	}
	return
}

func (w *Wallet) removeUtxo(u *message.Utxo) {
	var newUtxos []*message.Utxo
	for _, utxo := range w.Utxos {
		if u.Hash != utxo.Hash && u.N != utxo.N {
			newUtxos = append(newUtxos, utxo)
		}
	}
	w.Utxos = newUtxos
}

func (w *Wallet) CreateTxOuts(toAddress string, value, feeValue uint64) []*message.TxOut {
	var txout []*message.TxOut
	lockingScript1 := script.CreateLockingScriptForPKH(util.DecodeAddress(toAddress))
	txout = append(txout, &message.TxOut{
		Value:         value,
		LockingScript: common.NewVarStr(lockingScript1),
	})

	lockingScript2 := script.CreateLockingScriptForPKH(util.DecodeAddress(w.GetAddress()))
	txout = append(txout, &message.TxOut{
		Value:         feeValue,
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

		var sigHashCode = []byte{0x01, 0x00, 0x00, 0x00} // sig hash all
		sigbatureHash := util.Hash256(bytes.Join([][]byte{
			tx.Encode(),
			sigHashCode,
		}, []byte{}))

		signature, err := w.Sign(sigbatureHash)
		if err != nil {
			return nil, err
		}
		log.Printf("signature len: %+v", len(signature))
		var sigHashType = []byte{0x01}
		signatureWithType := bytes.Join([][]byte{signature, sigHashType}, []byte{})
		txin.UnlockingScript = script.CreateUnlockingScriptForPKH(signatureWithType, w.GetPublicKey())
		txins = append(txins, txin)
	}
	log.Printf("==== txins len: %+v", len(txins))
	return txins, nil
}

func (w *Wallet) Sign(sigHash []byte) ([]byte, error) {
	return w.Key.Sign(sigHash)
}
