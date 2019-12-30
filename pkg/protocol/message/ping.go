package message

import "encoding/binary"

type Ping struct {
	Nonce uint64
}

func DecodePing(b []byte) *Ping {
	return &Ping{
		Nonce: binary.LittleEndian.Uint64(b[0:8]),
	}
}
