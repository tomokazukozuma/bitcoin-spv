package client

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

type Client struct {
	Conn net.Conn
}

func NewClient(ip string) *Client {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{Conn: conn}
}

func (c *Client) SendMessage(msg message.Message) (int, error) {
	var checksum [4]byte
	hashedMsg := util.Hash256(msg.Encode())
	copy(checksum[:], hashedMsg[0:4])
	message := &common.Message{
		Magic:    binary.LittleEndian.Uint32([]byte{0x0B, 0x11, 0x09, 0x07}),
		Command:  msg.Command(),
		Length:   uint32(len(msg.Encode())),
		Checksum: checksum,
		Payload:  msg.Encode(),
	}
	return c.Conn.Write(message.Encode())
}

func (c *Client) ReceiveMessage(size uint32) ([]byte, error) {
	buf := make([]byte, size)
	_, err := c.Conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
