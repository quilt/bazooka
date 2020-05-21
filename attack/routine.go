package attack

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	SendTxs = iota
	SendBlock
	Sleep
	Exit
)

type Transaction struct {
	To       common.Address
	From     common.Address
	Nonce    uint64
	Data     hexutil.Bytes
	Amount   uint64
	GasPrice uint64 `yaml:"gas-price"`
	GasLimit uint64 `yaml:"gas-limit"`
}

func (tx *Transaction) toEthType() *types.Transaction {
	ret := types.NewTransaction(tx.Nonce, tx.To, big.NewInt(int64(tx.Amount)), tx.GasLimit, big.NewInt(int64(tx.GasPrice)), tx.Data)
	return ret
}

type Routine struct {
	Ty                 int
	Transactions       []*Transaction
	SignedTransactions types.Transactions
	SignedBlock        *types.Block
	SleepDuration      time.Duration `yaml:"sleep-duration"`
}
