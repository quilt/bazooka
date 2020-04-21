package main

import (
	"fmt"

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

	var (
		sentBlocks  = 0
		emptyBlocks = 0
	)

	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			return fmt.Errorf("failed to receive message from peer: %w", err)
		}

		switch {
		case msg.Code == eth.GetBlockHeadersMsg:
			if _, err = sp.handleGetBlockHeaderMsg(msg, rw, emptyBlocks); err != nil {
				return err
			}
		case msg.Code == eth.GetBlockBodiesMsg:
			if _, err = sp.handleGetBlockBodiesMsg(msg, rw, sentBlocks); err != nil {
				return err
			}
		case msg.Code == eth.NewBlockHashesMsg:
			if _, err = sp.handleNewBlockHashesMsg(msg, rw); err != nil {
				return err
			}
		default:
			log.Trace("Unrecognized message", "msg", msg)
		}
	}
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

func (sp *SimulationProtocol) handleGetBlockHeaderMsg(msg p2p.Msg, rw p2p.MsgReadWriter, emptyBlocks int) (int, error) {
	newEmptyBlocks := emptyBlocks
	var query eth.GetBlockHeadersData
	if err := msg.Decode(&query); err != nil {
		return newEmptyBlocks, fmt.Errorf("failed to decode msg %v: %w", msg, err)
	}
	log.Trace("GetBlockHeadersMsg", "query", query)
	headers := []*types.Header{}
	if query.Origin.Hash == (common.Hash{}) && !query.Reverse {
		number := query.Origin.Number
		for i := 0; i < int(query.Amount); i++ {
			if header := sp.chain.GetHeaderByNumber(number); header != nil {
				//fmt.Printf("Going to send block %d\n", header.Number.Uint64())
				headers = append(headers, header)
				if header.TxHash == types.EmptyRootHash {
					if !sp.markBlockSent(uint(number)) {
						newEmptyBlocks++
					}
				}
			} else {
				//fmt.Printf("Could not find header with number %d\n", number)
			}
			number += query.Skip + 1
		}
	}
	if query.Origin.Hash != (common.Hash{}) && query.Amount == 1 && query.Skip == 0 && !query.Reverse {
		if header := sp.chain.GetHeaderByHash(query.Origin.Hash); header != nil {
			log.Trace("Going to send header", "number", header.Number.Uint64())
			headers = append(headers, header)
		}
	}
	if err := p2p.Send(rw, eth.BlockHeadersMsg, headers); err != nil {
		return newEmptyBlocks, fmt.Errorf("failed to send headers: %w", err)
	}
	log.Info(fmt.Sprintf("Sent %d headers, empty blocks so far %d", len(headers), newEmptyBlocks))
	return newEmptyBlocks, nil
}

func (sp *SimulationProtocol) handleGetBlockBodiesMsg(msg p2p.Msg, rw p2p.MsgReadWriter, sentBlocks int) (int, error) {
	// Decode the retrieval message
	msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
	if _, err := msgStream.List(); err != nil {
		return 0, err
	}
	// Gather blocks until the fetch or network limits is reached
	var (
		hash   common.Hash
		bytes  int
		bodies []rlp.RawValue
	)
	// Retrieve the hash of the next block
	for {
		if err := msgStream.Decode(&hash); err == rlp.EOL {
			break
		} else if err != nil {
			return 0, fmt.Errorf("msg %v: %v", msg, err)
		}
		// Retrieve the requested block body, stopping if enough was found
		if data := sp.chain.GetBodyRLP(hash); len(data) != 0 {
			bodies = append(bodies, data)
			bytes += len(data)
		}
	}

	if err := p2p.Send(rw, eth.BlockBodiesMsg, bodies); err != nil {
		return len(bodies), err
	}

	log.Info("Sending bodies", "progress", len(bodies))

	return len(bodies), nil
}

func (sp *SimulationProtocol) sendLastBlock(rw p2p.MsgReadWriter) error {
	return p2p.Send(rw, eth.NewBlockMsg, []interface{}{sp.chain.CurrentBlock(), sp.chain.CurrentBlock().Difficulty()})
}

func (sp *SimulationProtocol) handleNewBlockHashesMsg(msg p2p.Msg, rw p2p.MsgReadWriter) (bool, error) {
	var blockHashMsg eth.NewBlockHashesData
	if err := msg.Decode(&blockHashMsg); err != nil {
		return false, fmt.Errorf("failed to decode msg %v: %w", msg, err)
	}
	log.Trace("NewBlockHashesMsg", "query", blockHashMsg)
	signaledHead := false
	for _, bh := range blockHashMsg {
		if bh.Number == sp.chain.CurrentBlock().NumberU64() {
			signaledHead = true
			break
		}
	}
	return signaledHead, nil
}
