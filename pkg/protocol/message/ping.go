package message

import (
	"encoding/binary"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol"
)

type Ping struct {
	Nonce uint64
}

func NewPing(nonce uint64) protocol.Message {
	return &Ping{
		Nonce: nonce,
	}
}
func (p *Ping) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "ping")
	return commandName
}

func (p *Ping) Encode() []byte {
	var nonce [8]byte
	binary.LittleEndian.PutUint64(nonce[:8], p.Nonce)
	return nonce[:]
}

func DecodePing(b []byte) *Ping {
	return &Ping{
		Nonce: binary.LittleEndian.Uint64(b[0:8]),
	}
}
