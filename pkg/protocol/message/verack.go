package message

import "github.com/tomokazukozuma/bitcoin-spv/pkg/protocol"

type Verack struct{}

func NewVerack() protocol.Message {
	return &Verack{}
}

func (v *Verack) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "verack")
	return commandName
}

func (v *Verack) Encode() []byte {
	return []byte{}
}
