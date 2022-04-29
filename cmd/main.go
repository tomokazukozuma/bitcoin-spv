package main

import (
	"log"
	"os"

	"github.com/tomokazukozuma/bitcoin-spv/internal/spv"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/network"
)

func main() {

	if len(os.Args) < 2 {
		//fmt.Println(usage)
		os.Exit(1)
	}

	// connect tcp
	c := network.NewClient("seed.tbtc.petertodd.org:18333")
	//c := network.NewClient("18.224.59.186:18333")
	defer c.Close()
	log.Printf("remote addr： %s", c.RemoteAddress().String())

	// handshake
	spv := spv.NewSPV(c)
	if err := spv.Handshake(0); err != nil {
		log.Fatal("handshake error: ", err)
	}
	log.Printf("address: %s", spv.Wallet.GetAddress())

	command := os.Args[1]
	switch command {
	case "balance":
		// send filterload
		if err := spv.SendFilterLoad(); err != nil {
			log.Fatal("filterload error: ", err)
		}

		// send getblocks
		if err := spv.SendGetBlocks("000000000000014ad045b835a2f4990a6acedccd95e8e3f42c0fe8caccba05a5"); err != nil {
			log.Fatal("GetBlocks error: ", err)
		}
		// receiving message
		if err := spv.MessageHandlerForBalance(); err != nil {
			log.Fatal("main: message handler err:", err)
		}

		balance := spv.Wallet.GetBalance()
		log.Printf("Balance: %d", balance)
	case "send":
		log.Printf("send")
		// send filterload
		if err := spv.SendFilterLoad(); err != nil {
			log.Fatal("filterload error: ", err)
		}

		// send getblocks
		if err := spv.SendGetBlocks("000000000000014ad045b835a2f4990a6acedccd95e8e3f42c0fe8caccba05a5"); err != nil {
			log.Fatal("GetBlocks error: ", err)
		}
		// receiving message
		if err := spv.MessageHandlerForBalance(); err != nil {
			log.Fatal("main: message handler err:", err)
		}
		balance := spv.Wallet.GetBalance()
		log.Printf("Balance: %d", balance)
		transaction := spv.SendTxInv("mgavKSS3hKCAyLKFhy5VHTYu5CMj8AAxQV", 1000)
		if err := spv.MessageHandlerForSend(transaction); err != nil {
			log.Fatal("main: message handler err:", err)
		}
	default:
		log.Printf("no command")
	}

	log.Printf("finish")
}
