package simulator

import (
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/lightclient/bazooka/attack"
	"github.com/lightclient/bazooka/simulator/contracts"
)

func InitBlockchain(db ethdb.Database, accounts map[common.Address]attack.Account) (*core.BlockChain, error) {
	n := 10

	genesis, err := Genesis()
	if err != nil {
		panic(fmt.Errorf("Could not create genesis: %s", err))
	}

	coinbaseKey, err := crypto.HexToECDSA("ad0f3019b6b8634c080b574f3d8a47ef975f0e4b9f63e82893e9a7bb59c2d609")
	if err != nil {
		return nil, err
	}

	txOpts := bind.NewKeyedTransactor(coinbaseKey)
	txOpts.GasPrice = big.NewInt(1)
	txOpts.GasLimit = 20 * params.TxGas
	txOpts.Nonce = big.NewInt(0)
	var nonce int64 = -1

	backend := &NoopBackend{db: db, genesis: genesis}

	var deployer *contracts.Deployer
	var deployerAddress common.Address
	var deploy = func(code []byte, salt []byte) *types.Transaction {
		var fixedSalt [32]byte
		copy(fixedSalt[:], salt[:])

		nonce++
		txOpts.Nonce.SetInt64(nonce)
		tx, err := deployer.Deploy(txOpts, code, fixedSalt)
		if err != nil {
			panic(err)
		}
		return tx
	}

	genesisBlock := genesis.MustCommit(db)

	engine := ethash.NewFaker()
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, engine, vm.Config{}, nil)
	blocks, _ := core.GenerateChain(params.TestChainConfig, genesisBlock, engine, db, n, func(i int, b *core.BlockGen) {
		b.SetCoinbase(crypto.PubkeyToAddress(coinbaseKey.PublicKey))
		b.SetExtra(common.BigToHash(big.NewInt(42)).Bytes())

		var tx *types.Transaction

		// deploy deployer contract
		if i == 1 {
			nonce++
			txOpts.Nonce.SetInt64(nonce)
			deployerAddress, tx, deployer, err = contracts.DeployDeployer(txOpts, backend)
			if err != nil {
				panic(err)
			}

			b.AddTx(tx)

			log.Info(fmt.Sprintf("Deployer address: %s", deployerAddress.Hex()))
		}

		//  initialize create2 contracts
		if i == 2 {
			for _, account := range accounts {
				if account.Code != nil {
					tx = deploy(account.Code, account.Salt)
					b.AddTx(tx)
				}
			}
		}
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
