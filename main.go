package main

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
)

func main() {
	setupLogger()

	blockchain, err := initBlockchain(rawdb.NewMemoryDatabase())
	if err != nil {
		panic(fmt.Errorf("Error initializing chain: %s", err))
	}

	server := makeP2PServer(blockchain)
	err = server.Start()
	if err != nil {
		panic("Error starting server")
	}

	err = addLocalPeer(server)
	if err != nil {
		panic(fmt.Errorf("Error adding local peer: %s", err))
	}

	time.Sleep(30 * time.Second)
	server.Stop()
}
