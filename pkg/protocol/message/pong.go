package message

import (
	"encoding/binary"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol"
)

type Pong struct {
	Nonce uint64
}

func NewPong(nonce uint64) protocol.Message {
	return &Pong{
		Nonce: nonce,
	}
}
func (p *Pong) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "pong")
	return commandName
}

func (p *Pong) Encode() []byte {
	var nonce [8]byte
	binary.LittleEndian.PutUint64(nonce[:8], p.Nonce)
	return nonce[:]
}

func DecodePong(b [8]byte) *Pong {
	return &Pong{
		Nonce: binary.LittleEndian.Uint64(b[0:8]),
	}
}
