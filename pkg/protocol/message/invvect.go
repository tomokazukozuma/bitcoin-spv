package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	InvTypeError = iota
	InvTypeMsgTx
	InvTypeMsgBlock
	InvTypeMsgFilteredBlock
	InvTypeMsgCmpctBlock
)

const InventoryVectorSize = 36

type InvVect struct {
	Type uint32
	Hash [32]byte
}

func NewInvVect(invType uint32, hash [32]byte) *InvVect {
	return &InvVect{
		Type: invType,
		Hash: hash,
	}
}

func DecodeInvVect(b []byte) (*InvVect, error) {
	if len(b) != InventoryVectorSize {
		return nil, fmt.Errorf("Decode to InvVect failed, invalid input: %v", b)
	}
	var arr [32]byte
	copy(arr[:], b[4:36])
	return &InvVect{
		Type: binary.LittleEndian.Uint32(b[0:4]),
		Hash: arr,
	}, nil
}

func (vect *InvVect) Encode() []byte {
	invType := make([]byte, 4)
	binary.LittleEndian.PutUint32(invType, vect.Type)
	return bytes.Join([][]byte{
		invType,
		vect.Hash[:],
	}, []byte{})
}
