package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("tra").(*LocalTransport)
	trb := NewLocalTransport("trb").(*LocalTransport)

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers["trb"], trb)
	assert.Equal(t, trb.peers["tra"], tra)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("tra")
	trb := NewLocalTransport("trb")

	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("hello message")

	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))

	rpc := <-trb.Consume()

	assert.Equal(t, rpc.Payload, Payload(rpc.Payload))
	assert.Equal(t, rpc.From, tra.Addr())
}
