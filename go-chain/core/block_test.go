package core

import (
	"bytes"
	"go-chain/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeaderEncodeDecode(t *testing.T) {
	header := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Timestamp: uint64(time.Now().Unix()),
		Height:    4,
		Nonce:     87654,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, header.EncodeBinary(buf))

	newHeader := &Header{}
	assert.Nil(t, newHeader.DecodeBinary(buf))
	assert.Equal(t, header, newHeader)
}

func TestBlockEncodeDecode(t *testing.T) {
	block := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: uint64(time.Now().Unix()),
			Height:    4,
			Nonce:     87654,
		},
		Transactions: nil,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, block.EncodeBinary(buf))

	newBlock := &Block{}
	assert.Nil(t, newBlock.DecodeBinary(buf))
	assert.Equal(t, block, newBlock)
}

func TestBlockHash(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: uint64(time.Now().Unix()),
			Height:    4,
			Nonce:     87654,
		},
		Transactions: []Transaction{},
	}

	hash := b.Hash()

	// t.Log(hash)

	assert.False(t, hash.IsZero())

}
