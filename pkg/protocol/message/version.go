package message

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

type Version struct {
	Version     uint32
	Services    uint64
	Timestamp   uint64
	AddrRecv    *common.NetworkAddress
	AddrFrom    *common.NetworkAddress
	Nonce       uint64
	UserAgent   *common.VarStr
	StartHeight uint32
	Relay       bool
}

func NewVersion() protocol.Message {
	addrFrom := &common.NetworkAddress{
		Services: uint64(1),
		IP: [16]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x7F, 0x00, 0x00, 0x01,
		},
		Port: 18333,
	}
	return &Version{
		Version:     uint32(70015),
		Services:    uint64(1),
		Timestamp:   uint64(time.Now().Unix()),
		AddrRecv:    addrFrom,
		AddrFrom:    addrFrom,
		Nonce:       uint64(0),
		UserAgent:   common.NewVarStr([]byte("")),
		StartHeight: uint32(0),
		Relay:       false,
	}
}
func (v *Version) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "version")
	return commandName
}

func (v *Version) Encode() []byte {
	var (
		version     [4]byte
		services    [8]byte
		timestamp   [8]byte
		addrRecv    [26]byte
		addrFrom    [26]byte
		nonce       [8]byte
		userAgent   []byte
		startHeight [4]byte
		relay       [1]byte
	)

	binary.LittleEndian.PutUint32(version[:4], v.Version)
	binary.LittleEndian.PutUint64(services[:8], v.Services)
	binary.LittleEndian.PutUint64(timestamp[:8], v.Timestamp)
	addrRecv = v.AddrRecv.Encode()
	addrFrom = v.AddrFrom.Encode()
	binary.LittleEndian.PutUint64(nonce[:8], v.Nonce)
	userAgent = v.UserAgent.Encode()
	binary.LittleEndian.PutUint32(startHeight[:4], v.StartHeight)
	if v.Relay {
		relay = [1]byte{0x01}
	} else {
		relay = [1]byte{0x00}
	}
	return bytes.Join(
		[][]byte{
			version[:],
			services[:],
			timestamp[:],
			addrRecv[:],
			addrFrom[:],
			nonce[:],
			userAgent[:],
			startHeight[:],
			relay[:],
		},
		[]byte{},
	)
}
