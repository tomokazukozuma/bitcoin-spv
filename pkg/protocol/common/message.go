package common

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
)

const (
	MagicTestnet3       = uint32(118034699)
	MessageHeaderLength = 24
)

type Message struct {
	Magic    uint32
	Command  [12]byte
	Length   uint32
	Checksum [4]byte
	Payload  []byte
}

func NewMessage(command [12]byte, payload []byte) *Message {
	var checksum [4]byte
	hashedPayload := util.Hash256(payload)
	copy(checksum[:], hashedPayload[0:4])
	return &Message{
		Magic:    MagicTestnet3,
		Command:  command,
		Length:   uint32(len(payload)),
		Checksum: checksum,
		Payload:  payload,
	}
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

func DecodeMessageHeader(b [MessageHeaderLength]byte) *Message {
	var (
		command  [12]byte
		checksum [4]byte
	)

	copy(command[:], b[4:16])
	copy(checksum[:], b[20:MessageHeaderLength])

	log.Printf("receive : %s", string(command[:]))

	return &Message{
		Magic:    binary.LittleEndian.Uint32(b[0:4]),
		Command:  command,
		Length:   binary.LittleEndian.Uint32(b[16:20]),
		Checksum: checksum,
	}
}

func IsValidChecksum(checksum [4]byte, payload []byte) bool {
	hashedPayload := util.Hash256(payload)
	var payloadChecksum [4]byte
	copy(payloadChecksum[:], hashedPayload[0:4])
	return checksum == payloadChecksum
}

func IsTestnet3(magic uint32) bool {
	return magic == MagicTestnet3
}
