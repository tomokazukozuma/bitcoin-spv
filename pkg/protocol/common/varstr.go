package common

import "bytes"

type VarStr struct {
	Length *VarInt
	Data   []byte
}

func NewVarStr(b []byte) *VarStr {
	len := uint64(len(b))
	length := NewVarInt(len)
	return &VarStr{
		Length: length,
		Data:   b,
	}
}

// Encode encode VarStr to byte slice.
func (s *VarStr) Encode() []byte {
	return bytes.Join([][]byte{
		s.Length.Encode(),
		s.Data,
	},
		[]byte{},
	)
}
