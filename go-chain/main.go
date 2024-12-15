package main

import (
	"fmt"
	"go-chain/network"
	"math/rand"
	"time"
)

func main() {
	localTr := network.NewLocalTransport("LOCAL")
	remoteTr := network.NewLocalTransport("REMOTE")
	// connect transports
	localTr.Connect(remoteTr)
	remoteTr.Connect(localTr)

	transports := []network.Transport{localTr, remoteTr}

	serverOpts := network.ServerOpts{
		Transports: transports,
	}
	server := network.NewServer(serverOpts)

	// send message every 3 seconds to local transport
	go func() {
		for {
			rand := rand.Intn(9000)

			remoteTr.SendMessage(localTr.Addr(), []byte(fmt.Sprintf("block number from remote: %d", rand)))
			time.Sleep(2 * time.Second)
		}

	}()

	// time.Sleep(1 * time.Second)

	go func() {
		for {
			rand := rand.Intn(9000)

			localTr.SendMessage(remoteTr.Addr(), []byte(fmt.Sprintf("block number from local: %d", rand)))
			time.Sleep(3 * time.Second)
		}

	}()

	fmt.Println("Starting go chain ...")
	server.Start()

}
