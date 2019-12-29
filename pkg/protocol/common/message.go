package common

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/util"
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

func IsValidChecksum(checksum [4]byte, payload []byte) bool {
	hashedPayload := util.Hash256(payload)
	var payloadChecksum [4]byte
	copy(payloadChecksum[:], hashedPayload[0:4])
	log.Printf("checksum: %+v", checksum)
	log.Printf("payloadChecksum: %+v", payloadChecksum)
	return checksum == payloadChecksum
}
