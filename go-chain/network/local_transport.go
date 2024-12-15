package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	peers     map[NetAddr]*LocalTransport
	lock      sync.RWMutex
	consumeCh chan RPC
}

func NewLocalTransport(addr string) Transport {
	return &LocalTransport{
		addr:      NetAddr(addr),
		peers:     make(map[NetAddr]*LocalTransport),
		consumeCh: make(chan RPC, 1024),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

func (t *LocalTransport) Connect(peer Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[peer.Addr()] = peer.(*LocalTransport)

	return nil
}

func (t *LocalTransport) SendMessage(addr NetAddr, msg Payload) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[addr]

	if !ok {
		return fmt.Errorf("%s: unable to send message to %s", t.addr, addr)
	}

	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: msg,
	}

	return nil
}
