package common

import "encoding/binary"

type NetworkAddress struct {
	Services uint64
	IP       [16]byte
	Port     uint16
}

func (addr *NetworkAddress) Encode() [26]byte {
	var b [26]byte
	binary.LittleEndian.PutUint64(b[0:8], addr.Services)
	copy(b[8:24], addr.IP[:])
	binary.BigEndian.PutUint16(b[24:26], addr.Port)
	return b
}
