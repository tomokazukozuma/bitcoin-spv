package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/internal/spv"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/client"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

func main() {

	// connect tcp
	c := client.NewClient("seed.tbtc.petertodd.org:18333")
	//c := client.NewClient("[2604:a880:400:d0::4ac1:9001]:18333")
	//c := client.NewClient("[2001:41d0:a:f7eb::1]:18333")
	//[2604:a880:2:d0::2065:5001]:18333 <-取得できたノード
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

	// handshake
	spv := spv.NewSPV(c)
	if err := spv.Handshake(); err != nil {
		log.Fatal("handshake error: ", err)
	}
	log.Printf("address: %s", spv.Wallet.GetAddress())

	// send filterload
	publicKeyHash := spv.Wallet.GetPublicKeyHash()
	spv.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{publicKeyHash}))

	// send getblocks
	startBlockHash, err := hex.DecodeString("00000000000000183a8addcaaa69faa69b5c818589df645212a349680c1af053")
	if err != nil {
		fmt.Println(err.Error())
	}
	//endBlockHash, err := hex.DecodeString("00000000000001920452f880f211635922a692c3ac23cdd79c961d5c7128541d")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	var reversedStartBlockHash [32]byte
	//var reversedEndBlockHash [32]byte
	copy(reversedStartBlockHash[:], util.ReverseBytes(startBlockHash))
	//copy(reversedEndBlockHash[:], util.ReverseBytes(endBlockHash))
	getblocks := message.NewGetBlocks(uint32(70015), [][32]byte{reversedStartBlockHash}, message.ZeroHash)
	spv.Client.SendMessage(getblocks)

	// receiving message
	if err := spv.MessageHandler(); err != nil {
		log.Printf("main: message handler err:", err)
	}

	log.Printf("finish")
}
