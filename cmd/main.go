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
	//c := network.NewClient("seed.tbtc.petertodd.org:18333")
	c := network.NewClient("18.224.59.186:18333")
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

	// handshake
	spv := spv.NewSPV(c)
	if err := spv.Handshake(0); err != nil {
		log.Fatal("handshake error: ", err)
	}
	log.Printf("address: %s", spv.Wallet.GetAddress())

	command := os.Args[1]
	switch command {
	case "balance":
		log.Printf("balance")
		// send filterload
		if err := spv.SendFilterLoad(); err != nil {
			log.Fatal("filterload error: ", err)
		}

		// send getblocks
		if err := spv.SendGetBlocks("0000000000167921f328c518bbf74919738dd44061a341d988e0505023995b14"); err != nil {
			log.Fatal("GetBlocks error: ", err)
		}
		// receiving message
		if err := spv.MessageHandler(); err != nil {
			log.Fatal("main: message handler err:", err)
		}
	case "send":
		log.Printf("balance")
		//if len(os.Args) != 5 {
		//	os.Exit(1)
		//}
		//addr := os.Args[2]
		//amount, err := strconv.Atoi(os.Args[3])
		//if err != nil {
		//	fmt.Printf("Invalid input amount %v\n", os.Args[3])
		//}
		//fee, err := strconv.Atoi(os.Args[4])
		//if err != nil {
		//	fmt.Printf("Invalid input amount %v\n", os.Args[4])
		//}
		//
		//sendBitcoin(addr, amount, fee)
	default:
		log.Printf("no command")
	}

	log.Printf("finish")
}
