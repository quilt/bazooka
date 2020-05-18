package attack

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/yaml.v2"
)

type Account struct {
	Key     hexutil.Bytes
	Balance uint64
	Code    hexutil.Bytes
	Salt    hexutil.Bytes
}

type Attack struct {
	Accounts map[common.Address]Account
	Routines []Routine
}

func (a *Attack) Load(s string) error {
	data, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, a); err != nil {
		return err
	}

	return nil
}

func (a *Attack) SignAll() {
	for i, r := range a.Routines {
		if r.Ty == SendTxs || r.Ty == SendBlock {
			txs := make([]*types.Transaction, 0)

			for _, tx := range r.Transactions {
				priv, err := crypto.ToECDSA(a.Accounts[tx.From].Key)

				signedTx, err := types.SignTx(tx.toEthType(), types.NewEIP155Signer(big.NewInt(1337)), priv)
				if err != nil {
					panic(err)
				}

				txs = append(txs, signedTx)
			}

			a.Routines[i].SignedTransactions = txs

			fmt.Printf("we got %d unsigned txs and %d signed txs\n", len(r.Transactions), len(a.Routines[i].SignedTransactions))
		}
	}
}

func (a *Attack) MakeBlocks(bc *core.BlockChain) {
	current := bc.CurrentHeader().Number

	log.Info("Generating routine blocks")

	for i, r := range a.Routines {
		if r.Ty == SendBlock {
			number := current.Uint64()
			parent := bc.GetBlockByNumber(number)
			coinbase := parent.Coinbase()
			timestamp := uint64(time.Now().Unix())

			gasLimit := core.CalcGasLimit(parent, 9223372036854775807, 9223372036854775807)

			header := &types.Header{
				ParentHash: parent.Hash(),
				Number:     big.NewInt(int64(number + 1)),
				GasLimit:   gasLimit,
				Extra:      []byte{},
				Time:       timestamp,
			}

			header.Coinbase = coinbase
			bc.Engine().Prepare(bc, header)

			statedb, err := bc.StateAt(parent.Root())
			if err != nil {
				panic("couldn't get state")
			}

			gasPool := new(core.GasPool).AddGas(header.GasLimit)
			txCount := 0
			var txs []*types.Transaction
			var receipts []*types.Receipt
			var blockFull = gasPool.Gas() < params.TxGas
			for _, tx := range r.SignedTransactions {
				if blockFull {
					break
				}

				statedb.Prepare(tx.Hash(), common.Hash{}, txCount)
				snap := statedb.Snapshot()

				receipt, err := core.ApplyTransaction(
					bc.Config(),
					bc,
					&coinbase,
					gasPool,
					statedb,
					header, tx, &header.GasUsed, *bc.GetVMConfig(),
				)
				if err != nil {
					statedb.RevertToSnapshot(snap)
					break
				}
				txs = append(txs, tx)
				receipts = append(receipts, receipt)
				txCount++
				if gasPool.Gas() < params.TxGas {
					blockFull = true
					break
				}
			}

			block, err := bc.Engine().FinalizeAndAssemble(bc, header, statedb, txs, []*types.Header{}, receipts)
			if err != nil {
				panic("couldn't finalize and assemble")
			}

			fmt.Printf("%x", block.Header())

			a.Routines[i].SignedBlock = block
			bc.SetHead(current.Uint64())
		}
	}

}
