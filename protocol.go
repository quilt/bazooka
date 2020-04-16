package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
)

func runProtocol(bc *core.BlockChain, peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	err := syncHandshake(bc, rw)
	if err != nil {
		return fmt.Errorf("Handshake failed: %s", err)
	}

	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			return fmt.Errorf("failed to receive message from peer: %w", err)
		}

		switch {
		case msg.Code == eth.GetBlockHeadersMsg:
			log.Info("Received eth.GetBlockHeadersMsg")
		case msg.Code == eth.GetBlockBodiesMsg:
			log.Info("Received eth.GetBlockBodiesMsg")
		case msg.Code == eth.NewBlockHashesMsg:
			log.Info("Received eth.NewBlockHashesMsg")
		case msg.Code == eth.StatusMsg:
			log.Info("Received eth.StatusMsg")
		default:
			log.Trace("Next message", "msg", msg)
		}
	}
}

func syncHandshake(bc *core.BlockChain, rw p2p.MsgReadWriter) error {
	type statusData struct {
		ProtocolVersion uint32
		NetworkID       uint64
		TD              *big.Int
		Head            common.Hash
		Genesis         common.Hash
		ForkID          forkid.ID
	}

	err := p2p.Send(rw, 0x00, &statusData{
		ProtocolVersion: 64,
		NetworkID:       1337,
		TD:              bc.CurrentBlock().Difficulty(),
		Head:            bc.CurrentHeader().Hash(),
		Genesis:         bc.Genesis().Hash(),
		ForkID:          forkid.NewID(bc),
	})
	if err != nil {
		return fmt.Errorf("failed to send status message to peer: %w", err)
	}

	return nil
}
