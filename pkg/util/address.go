package util

import (
	"bytes"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/bech32"
)

func EncodeAddress(publicKeyBytes []byte) string {
	bs := bytes.Join([][]byte{
		{0x6F},
		Hash160(publicKeyBytes),
	}, []byte{})

	checksum := Hash256(bs)[:4]
	return base58.Encode(bytes.Join([][]byte{bs, checksum}, []byte{}))
}

func DecodeAddress(address string) []byte {
	b := base58.Decode(address)
	pubkeyHashWithPrefix := b[:len(b)-4]
	checksum := Hash256(pubkeyHashWithPrefix)[:4]
	if !bytes.HasPrefix(checksum, b[len(b)-4:]) {
		log.Fatalf("mismatch checksum")
	}
	return pubkeyHashWithPrefix[1:] //1バイト外すかどうか
}

func EncodeNativeSegwitAddress(publicKeyBytes []byte) string {
	bs := bytes.Join([][]byte{
		[]byte{0x6F}, // This means that, this address is for testnet.
		Hash160(publicKeyBytes),
	},
		[]byte{})

	b5, err := bech32.ConvertBits(bs, 8, 5, true)
	log.Printf("err: %+v", err)
	b, _ := bech32.Encode("tb", b5)
	return b
}
