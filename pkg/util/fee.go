package util

import (
	"log"
	"math"
)

func CalculateFee(satoshiPerByte uint64, utxoCount int) uint64 {
	// baseTransactionSize = 8(Version + LockTime) + inputcounter + (txid + n + scriptLength + scriptsig(signature( 73) + pubkey(65)) + sequence)*utxoCount + outputcounter + output(8+1+24) *2
	var baseTransactionSize = 8 + 1 + (32+4+1+138+4)*utxoCount + 1 + 33*2

	totalTransactionSize := baseTransactionSize
	virtualTransactionSize := math.Ceil((float64(baseTransactionSize)*3 + float64(totalTransactionSize)) / 4)
	result := uint64(virtualTransactionSize) * satoshiPerByte
	log.Printf("tx size: %d", result)
	return result
}

func CalculateFeeForSegwit(satoshiPerByte uint64, utxoCount uint64) uint64 {
	// baseTransactionSize = 10(Version + Flag + Marker + LockTime) + inputcounter + (txid + n + scriptLength + scriptsig(<0 <20-byte-key-hash>>) + sequence)*utxoCount + outputcounter + output(8+1+23) *2
	var baseTransactionSize = 10 + 1 + (32+4+1+23+4)*utxoCount + 1 + 32*2

	//pushdata + signature(73,72,71) + pushdata + pubkey
	witnessSize := 1 + 73 + 1 + 33
	// baseTransactionSize + witnessCount +witnessSize*utxoCount
	totalTransactionSize := baseTransactionSize + 1 + uint64(witnessSize)*utxoCount
	virtualTransactionSize := math.Ceil((float64(baseTransactionSize)*3 + float64(totalTransactionSize)) / 4)
	return uint64(virtualTransactionSize) * satoshiPerByte
}
