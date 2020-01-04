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

func (w *SPV) MessageHandler() error {
	blockSize := 0
	needBlockSize := 1
	for {
		if needBlockSize == blockSize {
			log.Printf("====== break ======")
			return nil
		}
		buf, err := w.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
			return err
		}
		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)
		log.Printf("msg: %+v", msg)
		b, err := w.Client.ReceiveMessage(msg.Length)
		if err != nil {
			return err
		}
		if !common.IsTestnet3(msg.Magic) {
			log.Printf("not testnet3")
			continue
		}

		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendheaders")) {
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendcmpct")) {
		} else if bytes.HasPrefix(msg.Command[:], []byte("addr")) {
		} else if bytes.HasPrefix(msg.Command[:], []byte("getheaders")) {
		} else if bytes.HasPrefix(msg.Command[:], []byte("feefilter")) {
		} else if bytes.HasPrefix(msg.Command[:], []byte("ping")) {
			ping := message.DecodePing(b)
			pong := message.Pong{
				Nonce: ping.Nonce,
			}
			w.Client.SendMessage(&pong)
		} else if bytes.HasPrefix(msg.Command[:], []byte("inv")) {
			inv, _ := message.DecodeInv(b)
			log.Printf("inv.Count: %+v", inv.Count)

			inventory := []*message.InvVect{}
			for _, iv := range inv.Inventory {
				if iv.Type == message.InvTypeMsgBlock {
					inventory = append(inventory, message.NewInvVect(message.InvTypeMsgFilteredBlock, iv.Hash))
				}
			}
			log.Printf("inventory len: %+v", len(inventory))
			needBlockSize = len(inventory)
			log.Printf("needBlockSize: %+v", needBlockSize)
			_, err := w.Client.SendMessage(message.NewGetData(inventory))
			if err != nil {
				log.Fatalf("inv: send getdata message error: %+v", err)
			}
		} else if bytes.HasPrefix(msg.Command[:], []byte("merkleblock")) {
			if !common.IsValidChecksum(msg.Checksum, b) {
				log.Printf("invalid checksum")
				continue
			}
			blockSize++
			log.Printf("blockSize: %+v", blockSize)

			mb, _ := message.DecodeMerkleBlock(b)
			log.Printf("merkleblock: %+v", mb)
			log.Printf("hashCount: %+v", mb.HashCount.Data)
			if mb.HashCount.Data == 0 {
				continue
			}
			log.Printf("block hash: %s", mb.GetBlockHash())
			txHashes := mb.Validate()
			log.Printf("txHashes len: %+v", len(txHashes))
			for _, txHash := range txHashes {
				stringHash := hex.EncodeToString(util.ReverseBytes(txHash[:]))
				log.Printf("string txHash: %s", stringHash)
			}
			var inventory []*message.InvVect
			for _, txHash := range txHashes {
				inventory = append(inventory, message.NewInvVect(message.InvTypeMsgTx, txHash))
			}
			_, err := w.Client.SendMessage(message.NewGetData(inventory))
			if err != nil {
				log.Fatalf("merkleblock: send getdata message error: %+v", err)
			}
		} else if bytes.HasPrefix(msg.Command[:], []byte("tx")) {
			tx, _ := message.DecodeTx(b)
			log.Printf("tx: %+v", tx)
			log.Printf("txhash: %+v", tx.ID())
		} else if bytes.HasPrefix(msg.Command[:], []byte("notfound")) {
			getdata, _ := message.DecodeGetData(b)
			log.Printf("getdata: %+v", getdata)
			for _, v := range getdata.Inventory {
				log.Printf("inventory: %+v", v)
			}
		} else {
			log.Printf("receive : other")
		}
	}
}
