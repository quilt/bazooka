package payload

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

const (
	SendTxsRoutine = iota
	SendBlocksRoutine
	SleepRoutine
)

type Routine struct {
	ty            int
	from          int
	transactions  []types.Transaction
	blocks        []types.Block
	sleepDuration time.Duration
}
