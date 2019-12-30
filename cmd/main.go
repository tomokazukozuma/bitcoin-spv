package main

import (
	"bytes"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/internal/wallet"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/client"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
)

func main() {

	// connect tcp
	c := client.NewClient("testnet-seed.bitcoin.petertodd.org:18333")
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

	// handshake
	wallet := wallet.NewWallet(c)
	if err := wallet.Handshake(); err != nil {
		log.Fatal("handshake error: ", err)
	}

	// send filterload
	pubkey := bytes.Join([][]byte{wallet.Key.PublicKey.X.Bytes(), wallet.Key.PublicKey.Y.Bytes()}, []byte{})
	wallet.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{pubkey}))

	// receiving message
	wallet.MessageHandler()
	log.Printf("finish")
}
