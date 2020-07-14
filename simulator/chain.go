package simulator

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

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

type AccountWithAddress struct {
	addr    common.Address
	account attack.Account
}

func InitBlockchain(db ethdb.Database, height uint64, accountsMap map[common.Address]attack.Account) (*core.BlockChain, error) {
	accounts := make([]attack.Account, 0)
	accountsWithAddr := make([]AccountWithAddress, 0)

	sentBalances := false
	deployedAa := false
	for addr, account := range accountsMap {
		accounts = append(accounts, account)
		accountsWithAddr = append(accountsWithAddr, AccountWithAddress{addr: addr, account: account})
	}

	genesis, err := Genesis()
	if err != nil {
		return nil, err
	}

	coinbaseKey, err := crypto.HexToECDSA("ad0f3019b6b8634c080b574f3d8a47ef975f0e4b9f63e82893e9a7bb59c2d609")
	if err != nil {
		return nil, err
	}

	log.Info("coinbase", "addr", crypto.PubkeyToAddress(coinbaseKey.PublicKey))

	txOpts := bind.NewKeyedTransactor(coinbaseKey)
	txOpts.GasPrice = big.NewInt(1)
	txOpts.GasLimit = 20 * params.TxGas
	txOpts.Nonce = big.NewInt(0)
	var nonce uint64 = 0

	backend := &NoopBackend{db: db, genesis: genesis}

	var deployer *contracts.Deployer
	var deployerAddress common.Address
	var deploy = func(code []byte, salt []byte) *types.Transaction {
		var fixedSalt [32]byte
		copy(fixedSalt[:], salt[:])

		txOpts.Nonce.SetUint64(nonce)
		txOpts.GasLimit = 100000
		nonce++

		tx, err := deployer.Deploy(txOpts, code, fixedSalt)
		if err != nil {
			log.Error(fmt.Sprintf("Unable to deploy contract: %s", err))
			os.Exit(1)
		}
		return tx
	}

	var transfer = func(to common.Address, amt uint64) *types.Transaction {
		tx := types.NewTransaction(nonce, to, big.NewInt(int64(amt)), params.TxGas, nil, nil)
		tx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1337)), coinbaseKey)
		if err != nil {
			log.Error(fmt.Sprintf("Unable to sign tx: %s", err))
			os.Exit(1)
		}

		nonce++

		return tx
	}

	genesisBlock := genesis.MustCommit(db)

	engine := ethash.NewFaker()
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, engine, vm.Config{}, nil)
	blocks, _ := core.GenerateChain(genesis.Config, genesisBlock, engine, db, int(height), func(i int, b *core.BlockGen) {
		gasSpent := uint64(0)

		b.SetCoinbase(crypto.PubkeyToAddress(coinbaseKey.PublicKey))
		b.SetExtra(common.BigToHash(big.NewInt(42)).Bytes())

		var tx *types.Transaction

		// deploy deployer contract and initialize EOAs
		if i == 1 {
			txOpts.Nonce.SetUint64(nonce)
			deployerAddress, tx, deployer, err = contracts.DeployDeployer(txOpts, backend)

			if err != nil {
				log.Error(fmt.Sprintf("Unable to deploy Deployer: %s", err))
				os.Exit(1)
			}

			log.Info("create2 deployer deployed", "addr", deployerAddress)

			nonce++
			b.AddTx(tx)

		}

		// send balances
		if !sentBalances && i > 0 && uint64(i) < height {
			for len(accountsWithAddr) != 0 && gasSpent < 800000 {
				a := accountsWithAddr[0]
				accountsWithAddr = accountsWithAddr[1:]

				if a.account.Balance != 0 {
					tx := transfer(a.addr, a.account.Balance)
					b.AddTx(tx)
					gasSpent += tx.Gas()
				}
			}

			if len(accountsWithAddr) == 0 {
				sentBalances = true
			}
		}

		//  initialize create2 contracts
		if sentBalances && !deployedAa && i > 1 && uint64(i) < height {
			for len(accounts) != 0 && gasSpent < 800000 {
				account := accounts[0]
				accounts = accounts[1:]

				log.Info("contracts left", "amt", len(accounts), "spent", gasSpent, "code len", len(account.Code))

				// deploy AA
				if account.Code != nil {
					tx = deploy(account.Code, account.Salt)
					b.AddTx(tx)
					log.Info("deployed aa contact", "salt", account.Salt)
					gasSpent += tx.Gas()
					log.Info("deployed aa contact", "gas_spent", tx.Gas(), "salt", account.Salt)
				}
			}

			if len(accounts) == 0 {
				deployedAa = true
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
