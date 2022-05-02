package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReverseBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		exp   []byte
	}{
		{
			name:  "success",
			bytes: []byte{0x01, 0x02, 0x03, 0x04},
			exp:   []byte{0x04, 0x03, 0x02, 0x01},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b [4]byte
			copy(b[:], tt.bytes[:])
			ReverseBytes(b[:])
			assert.Equal(t, tt.exp, b[:])
		})
	}
}
