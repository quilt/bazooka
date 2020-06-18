package attack

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

func (a *Attack) SignAndAssemble(bc *core.BlockChain) error {
	// for now, assume the blocks will only be defined in a contiguous manner
	blockHeight := bc.CurrentHeader().Number.Uint64() + 1

	for i, r := range a.Routines {
		if r.Ty == SendTxs || r.Ty == SendBlock {
			txs, err := signTransactions(a.Accounts(), r.Transactions)
			if err != nil {
				return err
			}

			a.Routines[i].SignedTransactions = txs
		}

		if r.Ty == SendBlock {
			block, err := assembleBlock(bc, a.Routines[i].SignedTransactions, blockHeight)
			if err != nil {
				return err
			}

			blockHeight += 1
			a.Routines[i].SignedBlock = block
		}
	}

	bc.SetHead(a.Initialization.Height)

	return nil
}

func signTransactions(accounts map[common.Address]Account, txs []*Transaction) (types.Transactions, error) {
	stxs := make([]*types.Transaction, 0)

	for _, tx := range txs {
		// use key if exists
		if accounts[tx.From].Key.String() != "0x" {
			priv, err := crypto.ToECDSA(accounts[tx.From].Key)

			signedTx, err := types.SignTx(tx.toEthType(), types.NewEIP155Signer(big.NewInt(1337)), priv)
			if err != nil {
				return nil, err
			}

			stxs = append(stxs, signedTx)

		} else {
			// otherwise, assume it is AA
			stxs = append(stxs, tx.toEthType().WithAASignature())
		}
	}

	return stxs, nil
}

func assembleBlock(bc *core.BlockChain, txs []*types.Transaction, height uint64) (*types.Block, error) {
	parent := bc.GetBlockByNumber(height - 1)
	coinbase := parent.Coinbase()
	timestamp := uint64(parent.Time() + 10)
	gasLimit := core.CalcGasLimit(parent, 9223372036854775807, 9223372036854775807)

	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     big.NewInt(int64(height)),
		GasLimit:   gasLimit,
		Extra:      []byte{},
		Time:       timestamp,
		Coinbase:   coinbase,
	}

	bc.Engine().Prepare(bc, header)

	statedb, err := bc.StateAt(parent.Root())
	if err != nil {
		return nil, err
	}

	gasPool := new(core.GasPool).AddGas(header.GasLimit)
	var receipts []*types.Receipt
	var blockFull = gasPool.Gas() < params.TxGas

	for txCount, tx := range txs {
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
			log.Error("Transaction failed", "tx number", txCount)
			statedb.RevertToSnapshot(snap)
			break
		}

		receipts = append(receipts, receipt)

		if gasPool.Gas() < params.TxGas {
			blockFull = true
			break
		}
	}

	block, err := bc.Engine().FinalizeAndAssemble(bc, header, statedb, txs, []*types.Header{}, receipts)
	if err != nil {
		return nil, err
	}

	_, err = bc.InsertChain([]*types.Block{block})
	if err != nil {
		return nil, err
	}

	return block, nil
}
