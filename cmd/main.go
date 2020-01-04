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
	//c := client.NewClient("seed.tbtc.petertodd.org:18333")
	c := client.NewClient("[2001:41d0:a:f7eb::1]:18333")
	//[2604:a880:2:d0::2065:5001]:18333 <-取得できたノード
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

	// handshake
	spv := spv.NewSPV(c)
	if err := spv.Handshake(); err != nil {
		log.Fatal("handshake error: ", err)
	}

	// send filterload
	publicKeyHash := util.Hash160(spv.Key.PublicKey.SerializeUncompressed())
	spv.Client.SendMessage(message.NewFilterload(1024, 10, [][]byte{publicKeyHash}))

	// send getblocks
	startBlockHash, err := hex.DecodeString("000000000000020c54ca0a429835b14ba2f1629562547d39a0523af5dd518865")
	if err != nil {
		fmt.Println(err.Error())
	}
	var reversedStartBlockHash [32]byte
	copy(reversedStartBlockHash[:], util.ReverseBytes(startBlockHash))
	getblocks := message.NewGetBlocks(uint32(70015), [][32]byte{reversedStartBlockHash}, message.ZeroHash)
	spv.Client.SendMessage(getblocks)

	// receiving message
	if err := spv.MessageHandler(); err != nil {
		log.Printf("main: message handler err:", err)
	}

	log.Printf("finish")
}
