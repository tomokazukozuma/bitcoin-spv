package script

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/tomokazukozuma/bitcoin-spv/pkg/protocol/common"
)

func OpPushData(data []byte) []byte {
	len := len(data)
	if len <= 75 {
		return bytes.Join([][]byte{
			{byte(len)},
			data,
		}, []byte{})
	}
	if len <= math.MaxUint8 {
		return bytes.Join([][]byte{
			{OpPushData1},
			{byte(len)},
			data,
		}, []byte{})
	}
	if len <= math.MaxUint16 {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(len))
		return bytes.Join([][]byte{
			{OpPushData2},
			b,
			data,
		}, []byte{})
	}
	if len <= math.MaxUint32 {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(len))
		return bytes.Join([][]byte{
			{OpPushData4},
			b,
			data,
		}, []byte{})
	}
	return []byte{}
}

func CreateLockingScriptForPKH(pubkeyHash []byte) []byte {
	return bytes.Join([][]byte{
		{OpDup},
		{OpHash160},
		common.NewVarStr(pubkeyHash).Encode(),
		{OpEqualVerify},
		{OpCheckSig},
	}, []byte{})
}

func CreateUnlockingScriptForPKH(signature, publickey []byte) *common.VarStr {
	return common.NewVarStr(bytes.Join([][]byte{
		OpPushData(signature),
		OpPushData(publickey),
	}, []byte{}))
}
