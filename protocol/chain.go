package protocol

import (
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func InitBlockchain(db ethdb.Database) (*core.BlockChain, error) {
	n := 10

	genesis, err := Genesis()
	if err != nil {
		panic(fmt.Errorf("Could not create genesis: %s", err))
	}

	coinbaseKey, err := crypto.HexToECDSA("ad0f3019b6b8634c080b574f3d8a47ef975f0e4b9f63e82893e9a7bb59c2d609")
	if err != nil {
		return nil, err
	}
	coinbase := crypto.PubkeyToAddress(coinbaseKey.PublicKey)

	genesisBlock := genesis.MustCommit(db)

	engine := ethash.NewFaker()
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, engine, vm.Config{}, nil)
	blocks, _ := core.GenerateChain(params.TestChainConfig, genesisBlock, engine, db, n, func(i int, b *core.BlockGen) {
		b.SetCoinbase(coinbase)
		b.SetExtra(common.BigToHash(big.NewInt(42)).Bytes())
	})
	_, _ = blockchain.InsertChain(blocks)

	return blockchain, nil
}

func Genesis() (*core.Genesis, error) {
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
