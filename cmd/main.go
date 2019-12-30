package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"

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

	// send getblocks
	startBlockHash, err := hex.DecodeString("0000000000000657bda6681e1a3d1aac92d09d31721e8eedbca98cac73e93226")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var arr [32]byte
	copy(arr[:], util.ReverseBytes(startBlockHash))
	getblocks := message.NewGetBlocks(uint32(70015), [][32]byte{arr}, message.ZeroHash)
	wallet.Client.SendMessage(getblocks)

	// receiving message
	wallet.MessageHandler()
	log.Printf("finish")
}
