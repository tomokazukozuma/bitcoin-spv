package spv

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/network"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/wallet"
)

type SPV struct {
	Client *network.Client
	Wallet *wallet.Wallet
}

func NewSPV(client *network.Client) *SPV {
	wallet := wallet.NewWallet()
	return &SPV{
		Client: client,
		Wallet: wallet,
	}
}

func (s *SPV) Handshake() error {
	v := message.NewVersion()
	_, err := s.Client.SendMessage(v)
	if err != nil {
		return err
	}

	var recvVerack, sendVerack bool
	for {
		if recvVerack && sendVerack {
			log.Printf("success handshake")
			return nil
		}
		buf, err := s.Client.ReceiveMessage(common.MessageHeaderLength)
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
			_, err := s.Client.SendMessage(message.NewVerack())
			if err != nil {
				return err
			}
			sendVerack = true
		}
	}
}

func (s *SPV) SendFilterLoad() error {
	_, err := s.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{s.Wallet.GetPublicKeyHash()}))
	if err != nil {
		return err
	}
	return nil
}

func (s *SPV) SendGetBlocks(startBlockHeaderHash string) error {
	startBlockHash, err := hex.DecodeString(startBlockHeaderHash)
	if err != nil {
		return err
	}
	var reversedStartBlockHeaderHash [32]byte
	copy(reversedStartBlockHeaderHash[:], util.ReverseBytes(startBlockHash))
	getblocks := message.NewGetBlocks(uint32(70015), [][32]byte{reversedStartBlockHeaderHash}, message.ZeroHash)
	_, err = s.Client.SendMessage(getblocks)
	if err != nil {
		return err
	}
	return nil
}

func (s *SPV) MessageHandler() error {
	blockSize := 0
	needBlockSize := 1
	var transaction *message.Tx
	for {
		//if needBlockSize == blockSize {
		//	log.Printf("====== complete ======")
		//	return nil
		//}
		buf, err := s.Client.ReceiveMessage(common.MessageHeaderLength)
		if err != nil {
			log.Printf("ReceiveMessage: %+v", err)
			return err
		}
		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)
		log.Printf("msg: %+v", msg)
		b, err := s.Client.ReceiveMessage(msg.Length)
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
			pong := message.NewPong(ping.Nonce)
			s.Client.SendMessage(pong)
		} else if bytes.HasPrefix(msg.Command[:], []byte("inv")) {
			inv, _ := message.DecodeInv(b)
			log.Printf("inv.Count: %+v", inv.Count)

			inventory := []*common.InvVect{}
			for _, iv := range inv.Inventory {
				if iv.Type == common.InvTypeMsgBlock {
					inventory = append(inventory, common.NewInvVect(common.InvTypeMsgFilteredBlock, iv.Hash))
				}
			}
			log.Printf("inventory len: %+v", len(inventory))
			needBlockSize = len(inventory)
			log.Printf("needBlockSize: %+v", needBlockSize)
			_, err := s.Client.SendMessage(message.NewGetData(inventory))
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
			//log.Printf("merkleblock: %+v", mb)
			log.Printf("hashCount: %+v", mb.HashCount.Data)
			if mb.HashCount.Data == 0 {
				continue
			}
			log.Printf("block hash: %s", mb.GetBlockHash())
			txHashes := mb.Validate()
			//log.Printf("txHashes len: %+v", len(txHashes))
			for _, txHash := range txHashes {
				stringHash := hex.EncodeToString(util.ReverseBytes(txHash[:]))
				log.Printf("string txHash: %s", stringHash)
			}
			var inventory []*common.InvVect
			for _, txHash := range txHashes {
				inventory = append(inventory, common.NewInvVect(common.InvTypeMsgTx, txHash))
			}
			_, err := s.Client.SendMessage(message.NewGetData(inventory))
			if err != nil {
				log.Fatalf("merkleblock: send getdata message error: %+v", err)
			}
		} else if bytes.HasPrefix(msg.Command[:], []byte("tx")) {
			tx, _ := message.DecodeTx(b)
			log.Printf("tx: %+v", tx)
			log.Printf("txhash: %+v", tx.ID())
			utxos := tx.GetUtxo(s.Wallet.GetPublicKeyHash())
			for _, utxo := range utxos {
				s.Wallet.AddUtxo(utxo)
			}
			transaction = s.SendTxInv("mgavKSS3hKCAyLKFhy5VHTYu5CMj8AAxQV", 1000)
		} else if bytes.HasPrefix(msg.Command[:], []byte("getdata")) {
			getData, _ := message.DecodeGetData(b)
			log.Printf("getdata: %+v", getData)
			invs := getData.FilterInventoryWithType(common.InvTypeMsgTx)
			for _, invvect := range invs {
				txID := transaction.ID()
				if bytes.Equal(invvect.Hash[:], txID[:]) {
					fmt.Println("transaction send!")
					s.Client.SendMessage(transaction)
				}
			}
		} else if bytes.HasPrefix(msg.Command[:], []byte("notfound")) {
			getData, _ := message.DecodeGetData(b)
			log.Printf("getdata: %+v", getData)
			for _, v := range getData.Inventory {
				log.Printf("inventory: %+v", v)
			}
		} else {
			log.Printf("receive : other")
		}
	}
}

func (s *SPV) SendTxInv(toAddress string, value uint64) *message.Tx {
	transaction := s.Wallet.CreateTx(toAddress, value)
	inv := message.NewInv(
		common.NewVarInt(uint64(1)),
		[]*common.InvVect{common.NewInvVect(common.InvTypeMsgTx, transaction.ID())},
	).(*message.Inv)

	log.Printf("transaction: %+v", transaction)
	log.Printf("transaction txin count: %+v", transaction.TxInCount.Data)
	log.Printf("transaction txout count: %+v", transaction.TxOutCount.Data)
	log.Printf("transaction.ID: %+v", transaction.ID())
	log.Printf("transaction encode: %+v", hex.EncodeToString(transaction.Encode()))
	log.Printf("inv count: %+v", inv.Count)

	_, err := s.Client.SendMessage(inv)
	if err != nil {
		log.Fatalf("tx: send inv message error: %+v", err)
	}
	return transaction
}
