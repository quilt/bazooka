package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func initBlockchain(db ethdb.Database) (*core.BlockChain, error) {
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

	return blockchain, nil
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
