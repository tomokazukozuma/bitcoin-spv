package common

import "encoding/binary"

type VarInt struct {
	Data uint64
}

func NewVarInt(u uint64) *VarInt {
	return &VarInt{
		Data: u,
	}
}

func (v *VarInt) Encode() []byte {
	if v.Data < 0xfd {
		return []byte{byte(v.Data)}
	}
	if v.Data <= 0xffff {
		b := make([]byte, 3)
		b[0] = byte(0xfd)
		binary.LittleEndian.PutUint16(b[1:], uint16(v.Data))
		return b
	}
	if v.Data <= 0xffffffff {
		b := make([]byte, 5)
		b[0] = byte(0xfe)
		binary.LittleEndian.PutUint32(b[1:], uint32(v.Data))
		return b
	}
	if v.Data <= 0xffffffffffffffff {
		b := make([]byte, 9)
		b[0] = byte(0xff)
		binary.LittleEndian.PutUint64(b[1:], v.Data)
		return b
	}
	return []byte{byte(v.Data)}
}
