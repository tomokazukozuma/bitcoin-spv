package common

import (
	"bytes"
	"fmt"
)

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

func DecodeVarStr(b []byte) (*VarStr, error) {
	length, err := DecodeVarInt(b)
	if err != nil {
		return nil, err
	}
	varintLen := len(length.Encode())
	varstrLen := length.Data + uint64(varintLen)
	if uint64(len(b)) < varstrLen {
		return nil, fmt.Errorf("Decode varstr failed, invalid input: %v", b)
	}
	str := b[varintLen:varstrLen]
	return &VarStr{
		Length: length,
		Data:   str,
	}, nil
}
