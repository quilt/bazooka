package simulator

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

type NoopBackend struct {
	db      ethdb.Database
	genesis *core.Genesis
}

func (*NoopBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	panic("must not be called")
}
func (*NoopBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	panic("must not be called")
}

func (*NoopBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	panic("must not be called")
}
func (*NoopBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	panic("must not be called")
}
func (*NoopBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	panic("must not be called")
}
func (*NoopBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	panic("must not be called")
}
func (*NoopBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil // nothing to do
}

func (b *NoopBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, _, _, _ := rawdb.ReadReceipt(b.db, txHash, b.genesis.Config)
	return receipt, nil
}

func (*NoopBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	panic("must not be called")
}
func (*NoopBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	panic("must not be called")
}
