package util

import (
	"crypto/rand"
	"io/ioutil"

	"github.com/btcsuite/btcd/btcec"
)

const keypath = "./privatekey"

type Key struct {
	PrivateKey *btcec.PrivateKey
	PublicKey  *btcec.PublicKey
}

func NewKey() *Key {
	return &Key{}
}

func (k *Key) GenerateKey() error {
	if existsFile(keypath) {
		randomBytes, err := ioutil.ReadFile(keypath)
		if err != nil {
			return err
		}
		privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), randomBytes)
		k.PrivateKey = privateKey
		k.PublicKey = publicKey
	} else {
		randomBytes, err := generateRandom()
		if err != nil {
			return err
		}
		writeFile(keypath, randomBytes)
		privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), randomBytes)
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

func generateRandom() ([]byte, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
