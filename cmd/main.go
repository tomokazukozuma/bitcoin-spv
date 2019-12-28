package main

import (
	"bytes"
	"log"
	"time"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/client"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
)

func main() {
	c := client.NewClient("testnet-seed.bitcoin.jonasschnelli.ch:18333")
	defer c.Conn.Close()
	log.Printf("remote addr： %s", c.Conn.RemoteAddr().String())

	addrFrom := &common.NetworkAddress{
		Services: uint64(1),
		IP: [16]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x7F, 0x00, 0x00, 0x01,
		},
		Port: 8333,
	}
	v := &message.Version{
		Version:     uint32(70015),
		Services:    uint64(1),
		Timestamp:   uint64(time.Now().Unix()),
		AddrRecv:    addrFrom,
		AddrFrom:    addrFrom,
		Nonce:       uint64(0),
		UserAgent:   common.NewVarStr([]byte("")),
		StartHeight: uint32(0),
		Relay:       false,
	}
	_, err := c.SendMessage(v)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := c.ReceiveMessage(common.MessageLen)
	if err != nil {
		log.Fatal(err)
	}

	var header [24]byte
	copy(header[:], buf)
	msg := common.DecodeMessageHeader(header)
	log.Printf("receive: %s %d", msg.Command, msg.Length)
	payload, err := c.ReceiveMessage(msg.Length)
	if err != nil {
		log.Fatal(err)
	}

	if bytes.HasPrefix(msg.Command[:], []byte("verack")) {
		log.Printf("receive verack: %+v", payload)
	} else if bytes.HasPrefix(msg.Command[:], []byte("version")) {
		log.Printf("receive version: %+v", payload)
		_, err := c.SendMessage(&message.Verack{})
		if err != nil {
			log.Fatal(err)
		}
	}

	buf2, err := c.ReceiveMessage(common.MessageLen)
	if err != nil {
		log.Fatal(err)
	}

	var header2 [24]byte
	copy(header2[:], buf2)
	msg2 := common.DecodeMessageHeader(header2)
	payload2, err := c.ReceiveMessage(msg2.Length)
	if err != nil {
		log.Fatal(err)
	}

	if bytes.HasPrefix(msg2.Command[:], []byte("verack")) {
		log.Printf("receive verack: %+v", payload2)
	} else if bytes.HasPrefix(msg2.Command[:], []byte("version")) {
		log.Printf("receive version: %+v", payload2)
		c.SendMessage(&message.Verack{})
	}
}
