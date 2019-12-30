package message

import (
	"bytes"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

type GetData struct {
	Count     *common.VarInt
	Inventory []*common.InvVect
}

func NewGetData(inventory []*common.InvVect) *GetData {
	length := len(inventory)
	count := common.NewVarInt(uint64(length))
	return &GetData{
		Count:     count,
		Inventory: inventory,
	}
}

func (g *GetData) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "getdata")
	return commandName
}

func (g *GetData) Encode() []byte {
	inventoryBytes := [][]byte{}
	for _, invvect := range g.Inventory {
		inventoryBytes = append(inventoryBytes, invvect.Encode())
	}
	return bytes.Join([][]byte{
		g.Count.Encode(),
		bytes.Join(inventoryBytes, []byte{}),
	}, []byte{})
}
