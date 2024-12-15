package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Hash [32]uint8

func HashFromBytes(bytes []byte) Hash {
	if len(bytes) != 32 {
		msg := fmt.Sprintf("the length of the bytes should be 32, you gave %d", len(bytes))
		panic(msg)
	}
	return Hash(bytes)

}

func (h Hash) ToSlice() []byte {
	buf := make([]byte, 32)
	for i := 0; i < 32; i++ {
		buf[i] = h[i]
	}

	return buf
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func RandomBytes() []byte {
	token := make([]byte, 32)
	rand.Read(token)
	return token
}

func RandomHash() Hash {
	return HashFromBytes(RandomBytes())
}
