package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

type SimulationProtocol struct {
	chain        *core.BlockChain
	blockMarkers []uint64
}

func NewProtocolManager(bc *core.BlockChain) *SimulationProtocol {
	return &SimulationProtocol{chain: bc}
}

func (sp *SimulationProtocol) markBlockSent(blockNumber uint) bool {
	lengthNeeded := (blockNumber+63)/64 + 1

	if lengthNeeded > uint(len(sp.blockMarkers)) {
		sp.blockMarkers = append(sp.blockMarkers, make([]uint64, lengthNeeded-uint(len(sp.blockMarkers)))...)
	}

	bitMask := (uint64(1) << (blockNumber & 63))
	result := (sp.blockMarkers[blockNumber/64] & bitMask) != 0
	sp.blockMarkers[blockNumber/64] |= bitMask

	return result
}

func runProtocol(sp *SimulationProtocol, peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	err := syncHandshake(sp.chain, rw)
	if err != nil {
		return fmt.Errorf("Handshake failed: %s", err)
	}

	syncComplete := false

	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			return fmt.Errorf("failed to receive message from peer: %w", err)
		}

		switch {
		case msg.Code == eth.GetBlockHeadersMsg:
			if err = sp.handleGetBlockHeaderMsg(msg, rw); err != nil {
				return err
			}
		case msg.Code == eth.GetBlockBodiesMsg:
			if err = sp.handleGetBlockBodiesMsg(msg, rw); err != nil {
				return err
			}
		case msg.Code == eth.NewBlockHashesMsg:
			if syncComplete, err = sp.handleNewBlockHashesMsg(msg, rw); err != nil {
				return err
			}
		default:
			log.Trace("Unrecognized message", "msg", msg)
		}

		// break after sync is complete
		if syncComplete {
			break
		}
	}

	// send invalid tx
	time.Sleep(1 * time.Second)
	log.Info("Sending TransactionMsg now")
	if err := p2p.Send(rw, eth.TransactionMsg, []types.Transaction{*types.NewTransaction(0, common.BigToAddress(big.NewInt(42)), big.NewInt(1337), 1000000, big.NewInt(12), []byte{0})}); err != nil {
		return fmt.Errorf("couldn't announce new txs")
	}

	return nil
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

func (sp *SimulationProtocol) handleGetBlockHeaderMsg(msg p2p.Msg, rw p2p.MsgReadWriter) error {
	var query eth.GetBlockHeadersData
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("failed to decode msg %v: %w", msg, err)
	}

	log.Trace("GetBlockHeadersMsg", "query", query)

	if query.Reverse {
		return fmt.Errorf("reverse not supported")
	}

	headers := []*types.Header{}

	// if selecting via hash, convert to number
	if query.Origin.Hash != (common.Hash{}) {
		header := sp.chain.GetHeaderByHash(query.Origin.Hash)
		if header != nil {
			query.Origin.Hash = common.Hash{}
			query.Origin.Number = header.Number.Uint64()
		} else {
			return fmt.Errorf("Could not find header with hash %d\n", query.Origin.Hash)
		}
	}

	// find hashes via number
	number := query.Origin.Number
	for i := 0; i < int(query.Amount); i++ {
		if header := sp.chain.GetHeaderByNumber(number); header != nil {
			headers = append(headers, header)
		}
		number += query.Skip + 1
	}

	if err := p2p.Send(rw, eth.BlockHeadersMsg, headers); err != nil {
		return fmt.Errorf("failed to send headers: %w", err)
	}

	return nil
}

func (sp *SimulationProtocol) handleGetBlockBodiesMsg(msg p2p.Msg, rw p2p.MsgReadWriter) error {
	log.Trace("GetBlockBodiesMsg")

	msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
	if _, err := msgStream.List(); err != nil {
		return err
	}

	var (
		hash   common.Hash
		bytes  int
		bodies []rlp.RawValue
	)

	for {
		if err := msgStream.Decode(&hash); err == rlp.EOL {
			break
		} else if err != nil {
			return fmt.Errorf("msg %v: %v", msg, err)
		}

		if data := sp.chain.GetBodyRLP(hash); len(data) != 0 {
			bodies = append(bodies, data)
			bytes += len(data)
		}
	}

	if err := p2p.Send(rw, eth.BlockBodiesMsg, bodies); err != nil {
		return err
	}

	return nil
}

func (sp *SimulationProtocol) handleNewBlockHashesMsg(msg p2p.Msg, rw p2p.MsgReadWriter) (bool, error) {
	var blockHashMsg eth.NewBlockHashesData
	if err := msg.Decode(&blockHashMsg); err != nil {
		return false, fmt.Errorf("failed to decode msg %v: %w", msg, err)
	}

	log.Trace("NewBlockHashesMsg", "query", blockHashMsg)

	syncComplete := false
	for _, bh := range blockHashMsg {
		if bh.Number == sp.chain.CurrentBlock().NumberU64() {
			syncComplete = true
			break
		}
	}
	return syncComplete, nil
}
