package client

import (
	"log"
	"net"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/message"
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
	message := common.NewMessage(msg.Command(), msg.Encode())
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
