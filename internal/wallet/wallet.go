package wallet

import (
	"bytes"
	"log"
	"time"

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
	addrFrom := &common.NetworkAddress{
		Services: uint64(1),
		IP: [16]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x7F, 0x00, 0x00, 0x01,
		},
		Port: 8333,
	}
	v := &message.Version{
		Version:     uint32(70015),
		Services:    uint64(1),
		Timestamp:   uint64(time.Now().Unix()),
		AddrRecv:    addrFrom,
		AddrFrom:    addrFrom,
		Nonce:       uint64(0),
		UserAgent:   common.NewVarStr([]byte("")),
		StartHeight: uint32(0),
		Relay:       false,
	}
	_, err := w.Client.SendMessage(v)
	if err != nil {
		log.Fatal(err)
	}

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
