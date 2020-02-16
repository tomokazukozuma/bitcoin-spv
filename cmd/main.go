package main

import (
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/network"

	"github.com/tomokazukozuma/bitcoin-spv/internal/spv"
)

func main() {

	// connect tcp
	c := network.NewClient("seed.tbtc.petertodd.org:18333")
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

	// handshake
	spv := spv.NewSPV(c)
	if err := spv.Handshake(0); err != nil {
		log.Fatal("handshake error: ", err)
	}
	log.Printf("address: %s", spv.Wallet.GetAddress())

	// send filterload
	if err := spv.SendFilterLoad(); err != nil {
		log.Fatal("filterload error: ", err)
	}

	// send getblocks
	if err := spv.SendGetBlocks("0000000000000010708ca3fad77d86d01d3e6bcd79e38a787f160bce23417c21"); err != nil {
		log.Fatal("GetBlocks error: ", err)
	}
	// receiving message
	if err := spv.MessageHandler(); err != nil {
		log.Fatal("main: message handler err:", err)
	}

	log.Printf("finish")
}
