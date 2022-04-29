package network

import (
	"log"
	"net"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

type Client interface {
	SendMessage(msg protocol.Message) (int, error)
	ReceiveMessage(size uint32) ([]byte, error)
	RemoteAddress() net.Addr
	Close()
}
type client struct {
	net.Conn
}

func NewClient(address string) Client {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	return &client{conn}
}

func (c *client) SendMessage(msg protocol.Message) (int, error) {
	message := common.NewMessage(msg.Command(), msg.Encode())
	log.Printf("send    : %s", string(message.Command[:]))
	return c.Conn.Write(message.Encode())
}

func (c *client) ReceiveMessage(size uint32) ([]byte, error) {
	buf := make([]byte, size)
	_, err := c.Conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *client) RemoteAddress() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *client) Close()  {
	c.Conn.Close()
}
