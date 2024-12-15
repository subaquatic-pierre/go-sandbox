package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
}

type Server struct {
	rpcChan    chan RPC
	quitChan   chan struct{}
	transports []Transport
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		rpcChan:    make(chan RPC, 1024),
		quitChan:   make(chan struct{}, 1),
		transports: opts.Transports,
	}
}

func (s *Server) Start() error {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)

free:
	// main server loop
	for {
		select {
		// main server loop, listen for rpc messages
		case <-s.rpcChan:
			rpc := <-s.rpcChan
			fmt.Printf("%+s\n", rpc.Payload)
		case <-s.quitChan:
			goto free
		case <-ticker.C:
			// fmt.Println("server does things every 5 seconds")
		}
	}
}

func (s *Server) initTransports() error {
	// create go routine for each transport
	for _, tr := range s.transports {
		go func(tr Transport) {
			// send rpc to server rpc chanel
			for rpc := range tr.Consume() {
				s.rpcChan <- rpc
			}

		}(tr)

	}
	return nil
}
