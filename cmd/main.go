package main

import (
	"bytes"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/internal/wallet"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/client"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
)

func main() {

	c := client.NewClient("testnet-seed.bitcoin.petertodd.org:18333")
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

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
		log.Printf("receive command: %s", string(msg.Command[:]))
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
