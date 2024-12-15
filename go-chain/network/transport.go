package network

type NetAddr string
type Payload []byte

type RPC struct {
	From    NetAddr
	Payload Payload
}

type Transport interface {
	Addr() NetAddr
	Connect(Transport) error
	SendMessage(NetAddr, Payload) error
	Consume() <-chan RPC
}
