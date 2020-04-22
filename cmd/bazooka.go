package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func Execute() {
	setupLogger()

	db := rawdb.NewMemoryDatabase()

	blockchain, err := payload.initBlockchain(db)
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

func setupLogger() {
	var (
		ostream log.Handler
		glogger *log.GlogHandler
	)

	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)

	if usecolor {
		output = colorable.NewColorableStderr()
	}

	ostream = log.StreamHandler(output, log.TerminalFormat(usecolor))
	glogger = log.NewGlogHandler(ostream)
	log.Root().SetHandler(glogger)
	glogger.Verbosity(log.Lvl(5))
}
