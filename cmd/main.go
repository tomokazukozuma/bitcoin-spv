package main

import (
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
	c := client.NewClient("[2604:a880:2:d0::2065:5001]:18333")
	//c := client.NewClient("[2604:a880:400:d0::4ac1:9001]:18333")
	//[2604:a880:2:d0::2065:5001]:18333 <-取得できたノード
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

	// handshake
	wallet := wallet.NewWallet(c)
	if err := wallet.Handshake(); err != nil {
		log.Fatal("handshake error: ", err)
	}

	// send filterload
	publicKeyHash := util.Hash160(wallet.Key.PublicKey.SerializeUncompressed())
	wallet.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{publicKeyHash}))

	// send getblocks
	startBlockHash, err := hex.DecodeString("000000000000020c54ca0a429835b14ba2f1629562547d39a0523af5dd518865")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	hashStop := message.ZeroHash[:]
	var arr [32]byte
	copy(arr[:], util.ReverseBytes(startBlockHash))
	var arrHashStop [32]byte
	copy(arrHashStop[:], util.ReverseBytes(hashStop))
	getblocks := message.NewGetBlocks(uint32(70015), [][32]byte{arr}, arrHashStop)
	wallet.Client.SendMessage(getblocks)

	// receiving message
	wallet.MessageHandler()
	log.Printf("finish")
}
