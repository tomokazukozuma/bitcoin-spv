package spv

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/script"

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
	return &SPV{
		Client:  client,
		Key:     key,
		Address: util.EncodeAddress(key.PublicKey.SerializeUncompressed()),
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

func (s *SPV) MessageHandler() error {
	blockSize := 0
	needBlockSize := 1
	var transaction *message.Tx
	for {
		//if needBlockSize == blockSize {
		//	log.Printf("====== complete ======")
		//	return nil
		//}
		buf, err := s.Client.ReceiveMessage(common.MessageLen)
		if err != nil {
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
			pong := message.Pong{
				Nonce: ping.Nonce,
			}
			s.Client.SendMessage(&pong)
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
			var inventory []*message.InvVect
			for _, txHash := range txHashes {
				inventory = append(inventory, message.NewInvVect(message.InvTypeMsgTx, txHash))
			}
			_, err := s.Client.SendMessage(message.NewGetData(inventory))
			if err != nil {
				log.Fatalf("merkleblock: send getdata message error: %+v", err)
			}
		} else if bytes.HasPrefix(msg.Command[:], []byte("tx")) {
			tx, _ := message.DecodeTx(b)
			log.Printf("tx: %+v", tx)
			log.Printf("txhash: %+v", tx.ID())
			pubkeyHash := util.Hash160(s.Key.PublicKey.SerializeUncompressed())
			utxo := tx.GetUtxo(pubkeyHash)
			fee := util.CalculateFee(10, 1)
			chargeValue := utxo[0].TxOut.Value - 1000 - fee
			txout := createTxOut("2NCnDx5Zm6LgYerCjYe5TSQPeSdtsUdkmzn", s.Address, 1000, chargeValue)
			txin, err := createTxIn(utxo, txout, s.Key)
			if err != nil {
				log.Fatalf("createTxIn: %+v", err)
			}

			transaction = message.NewTx(uint32(1), txin, txout, uint32(0))
			inv := message.NewInv(
				common.NewVarInt(uint64(1)),
				[]*message.InvVect{message.NewInvVect(message.InvTypeMsgTx, transaction.ID())},
			)
			log.Printf("transaction: %+v", transaction)
			log.Printf("transaction txin count: %+v", transaction.TxInCount.Data)
			log.Printf("transaction txout count: %+v", transaction.TxOutCount.Data)
			log.Printf("transaction.ID: %+v", transaction.ID())
			log.Printf("transaction encode: %+v", hex.EncodeToString(transaction.Encode()))
			log.Printf("inv count: %+v", inv.Count)
			for _, iv := range inv.Inventory {
				log.Printf("inv type: %+v", iv.Type)
				log.Printf("inv hash: %+v", iv.Hash)
			}
			_, err = s.Client.SendMessage(inv)
			if err != nil {
				log.Fatalf("tx: send inv message error: %+v", err)
			}
		} else if bytes.HasPrefix(msg.Command[:], []byte("getdata")) {
			getData, _ := message.DecodeGetData(b)
			log.Printf("getdata: %+v", getData)
			invs := getData.FilterInventoryWithType(message.InvTypeMsgTx)
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

func createTxIn(utxos []*message.Utxo, txouts []*message.TxOut, key *util.Key) ([]*message.TxIn, error) {
	var txins []*message.TxIn
	for _, utxo := range utxos {
		txin := &message.TxIn{
			PreviousOutput: &message.OutPoint{
				Hash: utxo.Hash,
				N:    utxo.N,
			},
			UnlockingScript: utxo.TxOut.LockingScript,
			Sequence:        0xFFFFFFFF,
		}

		tx := message.NewTx(
			uint32(1),
			[]*message.TxIn{txin},
			txouts,
			uint32(0),
		)

		verified := util.Hash256(bytes.Join([][]byte{
			tx.Encode(),
			[]byte{0x01, 0x00, 0x00, 0x00},
		}, []byte{}))

		signature, err := key.Sign(verified)
		if err != nil {
			return nil, err
		}
		log.Printf("signature len: %+v", len(signature))
		hashType := []byte{0x01}
		signatureWithType := bytes.Join([][]byte{signature, hashType}, []byte{})
		txin.UnlockingScript = script.CreateUnlockingScriptForPKH(signatureWithType, key.PublicKey.SerializeUncompressed())
		txins = append(txins, txin)
	}
	log.Printf("==== txins len: %+v", len(txins))
	return txins, nil
}

func createTxOut(toAddress string, chargeAddress string, value, chargeValue uint64) []*message.TxOut {
	var txout []*message.TxOut
	lockingScript1 := script.CreateLockingScriptForPKH(util.DecodeAddress(toAddress))
	txout = append(txout, &message.TxOut{
		Value:         value,
		LockingScript: common.NewVarStr(lockingScript1),
	})

	lockingScript2 := script.CreateLockingScriptForPKH(util.DecodeAddress(chargeAddress))
	txout = append(txout, &message.TxOut{
		Value:         chargeValue,
		LockingScript: common.NewVarStr(lockingScript2),
	})
	return txout
}
