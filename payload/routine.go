package payload

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

const (
	SendTxsRoutine = iota
	SendBlockRoutine
	SleepRoutine
	Exit
)

type Routine struct {
	Ty            int
	From          int
	Transactions  []types.Transaction
	Block         types.Block
	SleepDuration time.Duration
}
