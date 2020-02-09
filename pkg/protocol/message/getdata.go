package message

import (
	"bytes"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol"
	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

type GetData struct {
	Count     *common.VarInt
	Inventory []*common.InvVect
}

func NewGetData(inventory []*common.InvVect) protocol.Message {
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

func DecodeGetData(b []byte) (*GetData, error) {
	count, _ := common.DecodeVarInt(b)
	b = b[len(count.Encode()):]
	var inventory []*common.InvVect
	for i := 0; uint64(i) < count.Data; i++ {
		iv, _ := common.DecodeInvVect(b[:36*(i+1)])
		inventory = append(inventory, iv)
	}
	return &GetData{
		Count:     count,
		Inventory: inventory,
	}, nil
}

func (g *GetData) FilterInventoryByType(typ uint32) []*common.InvVect {
	inventory := []*common.InvVect{}
	for _, invvect := range g.Inventory {
		if invvect.Type == typ {
			inventory = append(inventory, invvect)
		}
	}
	return inventory
}
