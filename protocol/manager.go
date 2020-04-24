package protocol

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/lightclient/bazooka/routine"
)

type Manager struct {
	chain        *core.BlockChain
	blockMarkers []uint64
	Routines     chan routine.Routine
}

func NewManager(bc *core.BlockChain) *Manager {
	return &Manager{chain: bc, Routines: make(chan routine.Routine, 10)}
}

func (pm *Manager) markBlockSent(blockNumber uint) bool {
	lengthNeeded := (blockNumber+63)/64 + 1

	if lengthNeeded > uint(len(pm.blockMarkers)) {
		pm.blockMarkers = append(pm.blockMarkers, make([]uint64, lengthNeeded-uint(len(pm.blockMarkers)))...)
	}

	bitMask := (uint64(1) << (blockNumber & 63))
	result := (pm.blockMarkers[blockNumber/64] & bitMask) != 0
	pm.blockMarkers[blockNumber/64] |= bitMask

	return result
}

func syncHandshake(bc *core.BlockChain, rw p2p.MsgReadWriter) error {
	status := eth.StatusData{
		ProtocolVersion: 64,
		NetworkID:       1337,
		TD:              bc.CurrentBlock().Difficulty(),
		Head:            bc.CurrentHeader().Hash(),
		Genesis:         bc.Genesis().Hash(),
		ForkID:          forkid.NewID(bc),
	}

	log.Debug(fmt.Sprintf("%#v,", status))

	err := p2p.Send(rw, 0x00, status)
	if err != nil {
		return fmt.Errorf("failed to send status message to peer: %w", err)
	}

	return nil
}
