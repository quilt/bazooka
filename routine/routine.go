package routine

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

const (
	SendTxs = iota
	SendBlock
	Sleep
	Exit
)

type Routine struct {
	Ty            int
	From          int
	Transactions  []types.Transaction
	Block         types.Block
	SleepDuration time.Duration
}
