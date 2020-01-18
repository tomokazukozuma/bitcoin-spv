package wallet

import (
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
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
	w.Utxos = append(w.Utxos, utxo)
}

func (w *Wallet) GetBalance() uint64 {
	var balance uint64
	for _, v := range w.Utxos {
		balance += v.TxOut.Value
	}
	return balance
}

func (w *Wallet) Sign(sigHash []byte) ([]byte, error) {
	return w.Key.Sign(sigHash)
}
