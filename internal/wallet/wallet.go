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
	//log.Printf("address: %s", address)
	//log.Printf("key: %+v, %+v", key.PrivateKey.X.Bytes(), key.PrivateKey.Y.Bytes())
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
		return err
	}

	var recvVerack, recvVersion bool
	for {
		if recvVerack && recvVersion {
			log.Printf("success handshake")
			return nil
		}
		buf, err := w.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
			log.Printf("handshake Receive message error: %+v", err)
			return err
		}

		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)
		_, err = w.Client.ReceiveMessage(msg.Length)
		if err != nil {
			return err
		}

		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
			recvVerack = true
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
			recvVersion = true
			_, err := w.Client.SendMessage(&message.Verack{})
			if err != nil {
				return err
			}
		}
	}
}

func (w *Wallet) MessageHandler() {
	for {
		buf, err := w.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
			//log.Printf("message handler err: %+v", err)
			log.Fatal("message handler err: ", err)
			//continue
		}
		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)

		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
			w.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
			w.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendheaders")) {
			w.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendcmpct")) {
			w.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("ping")) {
			b, _ := w.Client.ReceiveMessage(msg.Length)
			ping := message.DecodePing(b)
			pong := message.Pong{
				Nonce: ping.Nonce,
			}
			w.Client.SendMessage(&pong)
		} else if bytes.HasPrefix(msg.Command[:], []byte("addr")) {
			w.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("getheaders")) {
			w.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("feefilter")) {
			w.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("inv")) {
			log.Printf("msg: %+v", msg)
			b, _ := w.Client.ReceiveMessage(msg.Length)
			inv, _ := message.DecodeInv(b)
			log.Printf("inv.Count: %+v", inv.Count)

			inventory := []*common.InvVect{}
			for _, iv := range inv.Inventory {
				if iv.Type == common.InvTypeMsgBlock {
					inventory = append(inventory, common.NewInvVect(common.InvTypeMsgFilteredBlock, iv.Hash))
				} else {
					//inventory = append(inventory, iv)
				}
			}
			w.Client.SendMessage(message.NewGetData(inventory))
		} else if bytes.HasPrefix(msg.Command[:], []byte("merkleblock")) {
			b, _ := w.Client.ReceiveMessage(msg.Length)
			mb, _ := message.DecodeMerkleBlock(b)
			log.Printf("merkleblock: %+v", mb)
		} else {
			log.Printf("receive : other")
			w.Client.ReceiveMessage(msg.Length)
		}
	}
}

//func (w *Wallet) GetAddress() string {
//	return
//}
