package wallet

import (
	"bytes"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/client"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type Wallet struct {
	Client  *client.Client
	Key     *util.Key
	Address string
	Balance uint64
}

func NewWallet(client *client.Client) *Wallet {
	key := util.NewKey()
	key.GenerateKey()
	address := util.EncodeAddress(bytes.Join([][]byte{key.PublicKey.X.Bytes(), key.PublicKey.Y.Bytes()}, []byte{}))
	log.Printf("address: %s", address)
	return &Wallet{
		Client:  client,
		Key:     key,
		Address: address,
		Balance: 0,
	}
}

func (w *Wallet) Handshake() error {
	var recvVerack, recvVersion bool
	for {
		if recvVerack && recvVersion {
			return nil
		}
		buf, err := w.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
			return err
		}

		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)
		payload, err := w.Client.ReceiveMessage(msg.Length)
		if err != nil {
			return err
		}

		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
			log.Printf("receive verack: %+v", payload)
			recvVerack = true
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
			recvVersion = true
			log.Printf("receive version: %+v", payload)
			_, err := w.Client.SendMessage(&message.Verack{})
			if err != nil {
				return err
			}
		}
	}
}

//func (w *Wallet) GetAddress() string {
//	return
//}
