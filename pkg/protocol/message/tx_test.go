package message

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DecodeTx_Encode(t *testing.T) {
	rawtx, _ := hex.DecodeString("02000000012caa664f48f3631658e0588815c04b2c22f02d5af545506aa01e93dee90d0f4e0000000000feffffff02f8c8455300000000160014ecb7a596cc48b95a10bdb622fa75d87795c06b0b10270000000000001976a9142c26d5493277afcf2185b8901e0e0b26f282699e88ac29e62100")
	tx, _ := DecodeTx(rawtx)
	assert.Equal(t, rawtx, tx.Encode())
}
