package common

import (
	"bytes"
	"encoding/binary"
)

const MessageLen = 24

type Message struct {
	Magic    uint32
	Command  [12]byte
	Length   uint32
	Checksum [4]byte
	Payload  []byte
}

func (m *Message) Encode() []byte {
	var (
		magic  [4]byte
		length [4]byte
	)
	binary.LittleEndian.PutUint32(magic[:], m.Magic)
	binary.LittleEndian.PutUint32(length[:], m.Length)
	return bytes.Join([][]byte{
		magic[:],
		m.Command[:],
		length[:],
		m.Checksum[:],
		m.Payload,
	},
		[]byte{},
	)
}

func DecodeMessageHeader(b [MessageLen]byte) *Message {
	var (
		command  [12]byte
		checksum [4]byte
	)

	copy(command[:], b[4:16])
	copy(checksum[:], b[20:MessageLen])
	return &Message{
		Magic:    binary.LittleEndian.Uint32(b[0:4]),
		Command:  command,
		Length:   binary.LittleEndian.Uint32(b[16:20]),
		Checksum: checksum,
	}
}
