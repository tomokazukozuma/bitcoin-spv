package main

import (
	"bytes"
	"log"
	"time"

	"github.com/tomokazukozuma/bitcoin-spv/internal/wallet"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/client"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
)

func main() {

	c := client.NewClient("testnet-seed.bitcoin.petertodd.org:18333")
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

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
	_, err := c.SendMessage(v)
	if err != nil {
		log.Fatal(err)
	}

	wallet := wallet.NewWallet(c)
	if err := wallet.Handshake(); err != nil {
		log.Fatal(err)
	}

	pubkey := bytes.Join([][]byte{wallet.Key.PublicKey.X.Bytes(), wallet.Key.PublicKey.Y.Bytes()}, []byte{})
	wallet.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{pubkey}))

	// tcp messageのハンドラーを実装
	// 必要ないメッセージは読み込んで全て捨てる処理を入れる必要あり
	size := uint32(common.MessageLen)
	for {
		buf, _ := wallet.Client.ReceiveMessage(size)
		var header [24]byte
		copy(header[:], buf)
		msg := common.DecodeMessageHeader(header)
		//log.Printf("receive command: %s", string(msg.Command[:]))
		log.Printf("receive msg.Length: %+v", msg.Length)
		if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
			log.Printf("receive verack")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
			log.Printf("receive version")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendheaders")) {
			log.Printf("receive sendheaders")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("sendcmpct")) {
			log.Printf("receive sendcmpct")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("ping")) {
			log.Printf("receive ping")
			b, _ := wallet.Client.ReceiveMessage(msg.Length)
			ping := message.DecodePing(b)
			pong := message.Pong{
				Nonce: ping.Nonce,
			}
			wallet.Client.SendMessage(&pong)
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("addr")) {
			log.Printf("receive addr")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("getheaders")) {
			log.Printf("receive getheaders")
			wallet.Client.ReceiveMessage(msg.Length)
		} else if bytes.HasPrefix(msg.Command[:], []byte("feefilter")) {
			log.Printf("receive feefilter")
			wallet.Client.ReceiveMessage(msg.Length)
		} else {
			log.Printf("receive other")
			wallet.Client.ReceiveMessage(msg.Length)
		}
	}
	log.Printf("finish")
}
