package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func main() {
	setupLogger()

	node, err := getTargetAddr()
	if err != nil {
		panic(fmt.Errorf("Couldn't get target enode: %s", err))
	}

	//
	db := rawdb.NewMemoryDatabase()
	genesis, err := genesis()
	if err != nil {
		panic(fmt.Errorf("Could not create genesis: %s", err))
	}
	genesisBlock := genesis.MustCommit(db)

	engine := ethash.NewFaker()
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, engine, vm.Config{}, nil)
	blocks, _ := core.GenerateChain(params.TestChainConfig, genesisBlock, engine, db, 1, func(i int, b *core.BlockGen) {
		b.SetCoinbase(common.Address{0: byte(1), 19: byte(i)})
	})
	_, _ = blockchain.InsertChain(blocks)
	//

	server := makeP2PServer(blockchain)
	err = server.Start()
	if err != nil {
		panic("Error starting server")
	}
	server.AddPeer(node)
	time.Sleep(10 * time.Second)
	server.Stop()
}

func makeP2PServer(bc *core.BlockChain) *p2p.Server {
	serverKey, err := crypto.GenerateKey()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate server key: %v", err))
	}

	p2pConfig := p2p.Config{
		PrivateKey: serverKey,
		Name:       "bazooka",
		Logger:     log.New(),
		MaxPeers:   1,
		Protocols: []p2p.Protocol{
			p2p.Protocol{
				Name:    "eth",
				Version: 64,
				Length:  17,
				Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
					return protocolRun(bc, peer, rw)
				},
			},
		},
	}

	return &p2p.Server{Config: p2pConfig}
}

func protocolRun(bc *core.BlockChain, peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	err := syncHandshake(bc, rw)
	if err != nil {
		return fmt.Errorf("Handshake failed: %s", err)
	}

	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			return fmt.Errorf("failed to receive message from peer: %w", err)
		}

		switch {
		case msg.Code == eth.GetBlockHeadersMsg:
			log.Info("Received eth.GetBlockHeadersMsg")
		case msg.Code == eth.GetBlockBodiesMsg:
			log.Info("Received eth.GetBlockBodiesMsg")
		case msg.Code == eth.NewBlockHashesMsg:
			log.Info("Received eth.NewBlockHashesMsg")
		case msg.Code == eth.StatusMsg:
			log.Info("Received eth.StatusMsg")
		default:
			log.Trace("Next message", "msg", msg)
		}
	}
}

func syncHandshake(bc *core.BlockChain, rw p2p.MsgReadWriter) error {

	type statusData struct {
		ProtocolVersion uint32
		NetworkID       uint64
		TD              *big.Int
		Head            common.Hash
		Genesis         common.Hash
		ForkID          forkid.ID
	}

	err := p2p.Send(rw, 0x00, &statusData{
		ProtocolVersion: 64,
		NetworkID:       1337,
		TD:              bc.CurrentBlock().Difficulty(),
		Head:            bc.CurrentHeader().Hash(),
		Genesis:         bc.Genesis().Hash(),
		ForkID:          forkid.NewID(bc),
	})
	if err != nil {
		return fmt.Errorf("failed to send status message to peer: %w", err)
	}

	return nil
}

func genesis() (*core.Genesis, error) {
	raw, err := ioutil.ReadFile("genesis.json")
	if err != nil {
		return nil, err
	}

	genesis := new(core.Genesis)

	err = genesis.UnmarshalJSON(raw)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}

func getTargetAddr() (*enode.Node, error) {
	nodeKeyHex, err := ioutil.ReadFile("/home/matt/eth/aasim/geth/nodekey")
	if err != nil {
		return nil, err
	}

	nodeKey, err := crypto.HexToECDSA(string(nodeKeyHex))
	nodeid := fmt.Sprintf("%x", crypto.FromECDSAPub(&nodeKey.PublicKey)[1:])

	addr := fmt.Sprint("enode://", nodeid, "@127.0.0.1:30303?discport=0")
	fmt.Println(addr)

	nodeToConnect, err := enode.ParseV4(string(addr))
	if err != nil {
		return nil, fmt.Errorf("could not parse the node info: %w", err)
	}

	log.Info("Parsed node: %s, IP: %s\n", nodeToConnect, nodeToConnect.IP())

	return nodeToConnect, nil
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
