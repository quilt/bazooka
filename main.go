package main

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
)

func main() {
	setupLogger()

	db := rawdb.NewMemoryDatabase()

	blockchain, err := initBlockchain(db)
	if err != nil {
		panic(fmt.Errorf("Error initializing chain: %s", err))
	}

	pw := NewProtocolManager(blockchain)

	server := makeP2PServer(pw)
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
