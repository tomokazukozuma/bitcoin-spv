package wallet

import (
	"bytes"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/client"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
)

type Wallet struct {
	Client  *client.Client
	Balance uint64
}

func NewWallet(client *client.Client) *Wallet {
	return &Wallet{
		Client:  client,
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
