package handler

import (
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func InitBlockchain(db ethdb.Database) (*core.BlockChain, error) {
	n := 10

	genesis, err := genesis()
	if err != nil {
		panic(fmt.Errorf("Could not create genesis: %s", err))
	}

	genesisBlock := genesis.MustCommit(db)

	engine := ethash.NewFaker()
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, engine, vm.Config{}, nil)
	blocks, _ := core.GenerateChain(params.TestChainConfig, genesisBlock, engine, db, n, func(i int, b *core.BlockGen) {
		b.SetCoinbase(common.BigToAddress(big.NewInt(1337)))
		b.SetExtra(common.BigToHash(big.NewInt(42)).Bytes())
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
