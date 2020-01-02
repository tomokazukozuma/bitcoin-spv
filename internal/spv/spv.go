package spv

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/client"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type SPV struct {
	Client  *client.Client
	Key     *util.Key
	Address string
	Balance uint64
}

func NewSPV(client *client.Client) *SPV {
	key := util.NewKey()
	key.GenerateKey()
	serializedPubKey := key.PublicKey.SerializeUncompressed()
	address := util.EncodeAddress(serializedPubKey)
	return &SPV{
		Client:  client,
		Key:     key,
		Address: address,
		Balance: 0,
	}
}

func (s *SPV) Handshake() error {
	v := message.NewVersion()
	_, err := s.Client.SendMessage(v)
	if err != nil {
		return err
	}

	var recvVerack, recvVersion bool
	for {
		if recvVerack && recvVersion {
			log.Printf("success handshake")
			return nil
		}
		buf, err := s.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
			log.Printf("handshake Receive message error: %+v", err)
			return err
		}

		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)
		_, err = s.Client.ReceiveMessage(msg.Length)
		if err != nil {
			return err
		}

		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
			recvVerack = true
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
			recvVersion = true
			_, err := s.Client.SendMessage(&message.Verack{})
			if err != nil {
				return err
			}
		}
	}
}

func (w *SPV) MessageHandler() {
	for {
		buf, err := w.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
			log.Fatal("message handler err: ", err)
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

			inventory := []*message.InvVect{}
			for _, iv := range inv.Inventory {
				if iv.Type == message.InvTypeMsgBlock {
					inventory = append(inventory, message.NewInvVect(message.InvTypeMsgFilteredBlock, iv.Hash))
				}
			}
			w.Client.SendMessage(message.NewGetData(inventory))
		} else if bytes.HasPrefix(msg.Command[:], []byte("merkleblock")) {
			b, _ := w.Client.ReceiveMessage(msg.Length)
			mb, _ := message.DecodeMerkleBlock(b)
			log.Printf("merkleblock: %+v", mb)
			h := mb.GetBlockHash()
			hexHash := hex.EncodeToString(util.ReverseBytes(h[:]))
			log.Printf("BlockHash: %s", hexHash)
			log.Printf("Hashes length: %+v", len(mb.Hashes))
			txHashes := mb.Validate()
			log.Printf("txHashes: %+v", txHashes)
			inventory := []*message.InvVect{}
			for _, txHash := range txHashes {
				inventory = append(inventory, message.NewInvVect(message.InvTypeMsgTx, txHash))
			}
			w.Client.SendMessage(message.NewGetData(inventory))
		} else if bytes.HasPrefix(msg.Command[:], []byte("tx")) {
			b, _ := w.Client.ReceiveMessage(msg.Length)
			tx, _ := message.DecodeTx(b)
			log.Printf("tx: %+v", tx)
			log.Printf("TxID: %+v", tx.ID())
		} else {
			log.Printf("receive : other")
			w.Client.ReceiveMessage(msg.Length)
		}
	}
}
