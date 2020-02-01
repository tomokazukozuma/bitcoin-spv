package util

import (
	"crypto/rand"
	"io/ioutil"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

const keypath = "./privatekey"

type Key struct {
	PrivateKey *btcec.PrivateKey
	PublicKey  *btcec.PublicKey
}

func NewKey() *Key {
	key := &Key{}
	key.GenerateKey()
	return key
}

func (k *Key) GenerateKey() error {
	if existsFile(keypath) {
		privKeyBytes, err := ioutil.ReadFile(keypath)
		if err != nil {
			return err
		}
		privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
		k.PrivateKey = privateKey
		k.PublicKey = publicKey
	} else {
		privKey, err := generatePrivKey()
		if err != nil {
			return err
		}
		writeFile(keypath, privKey.Bytes())
		privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), privKey.Bytes())
		k.PrivateKey = privateKey
		k.PublicKey = publicKey
	}
	return nil
}

func (k *Key) Sign(message []byte) ([]byte, error) {
	signature, err := k.PrivateKey.Sign(message)
	if err != nil {
		return nil, err
	}
	return signature.Serialize(), nil
}

func generatePrivKey() (*big.Int, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	privKey := new(big.Int).SetBytes(b)
	var one = new(big.Int).SetInt64(1)

	// 1 < privkey < (n-1) の範囲になるように調整
	n := new(big.Int).Sub(btcec.S256().N, one)
	privKey.Mod(privKey, n)
	privKey.Add(privKey, one)
	return privKey, nil
}
