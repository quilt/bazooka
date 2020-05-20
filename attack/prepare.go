package attack

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

func (a *Attack) SignAndAssemble(bc *core.BlockChain) error {
	for i, r := range a.Routines {
		if r.Ty == SendTxs || r.Ty == SendBlock {
			txs, err := signTransactions(a.Accounts, r.Transactions)
			if err != nil {
				return err
			}

			a.Routines[i].SignedTransactions = txs
		}

		if r.Ty == SendBlock {
			block, err := assembleBlock(bc, r.SignedTransactions)
			if err != nil {
				return err
			}

			a.Routines[i].SignedBlock = block
		}
	}

	return nil
}

func signTransactions(accounts map[common.Address]Account, txs []*Transaction) (types.Transactions, error) {
	stxs := make([]*types.Transaction, 0)

	for _, tx := range txs {
		priv, err := crypto.ToECDSA(accounts[tx.From].Key)

		signedTx, err := types.SignTx(tx.toEthType(), types.NewEIP155Signer(big.NewInt(1337)), priv)
		if err != nil {
			return nil, err
		}

		stxs = append(stxs, signedTx)
	}

	return stxs, nil
}

func assembleBlock(bc *core.BlockChain, txs []*types.Transaction) (*types.Block, error) {
	current := bc.CurrentHeader().Number.Uint64()
	parent := bc.GetBlockByNumber(current)
	coinbase := parent.Coinbase()
	timestamp := uint64(time.Now().Unix())

	gasLimit := core.CalcGasLimit(parent, 9223372036854775807, 9223372036854775807)

	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     big.NewInt(int64(current + 1)),
		GasLimit:   gasLimit,
		Extra:      []byte{},
		Time:       timestamp,
	}

	header.Coinbase = coinbase
	bc.Engine().Prepare(bc, header)

	statedb, err := bc.StateAt(parent.Root())
	if err != nil {
		return nil, err
	}

	gasPool := new(core.GasPool).AddGas(header.GasLimit)
	txCount := 0
	var receipts []*types.Receipt
	var blockFull = gasPool.Gas() < params.TxGas

	for _, tx := range txs {
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
		receipts = append(receipts, receipt)
		txCount++
		if gasPool.Gas() < params.TxGas {
			blockFull = true
			break
		}
	}

	block, err := bc.Engine().FinalizeAndAssemble(bc, header, statedb, txs, []*types.Header{}, receipts)
	if err != nil {
		return nil, err
	}

	bc.SetHead(current)

	return block, nil
}
