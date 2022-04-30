package main

import (
	"log"
	"os"

	"github.com/tomokazukozuma/bitcoin-spv/internal/spv"
)

func main() {

	if len(os.Args) < 2 {
		//fmt.Println(usage)
		os.Exit(1)
	}

	spv := spv.NewSPV()
	defer spv.Close()
	if err := spv.Handshake(0); err != nil {
		log.Fatal("handshake error: ", err)
	}
	log.Printf("address: %s", spv.GetAddress())

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

		balance := spv.GetBalance()
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
		balance := spv.GetBalance()
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
