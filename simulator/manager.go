package simulator

import (
	"github.com/ethereum/go-ethereum/core"
	ep2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/lightclient/bazooka/attack"
	"github.com/lightclient/bazooka/p2p"
	"github.com/lightclient/bazooka/protocol"
)

type Manager struct {
	pms           []*protocol.Manager
	servers       []*ep2p.Server
	attack        attack.Runner
	targetDataDir string
}

func NewManager(chain *core.BlockChain, targetDataDir string) Manager {
	var pms []*protocol.Manager

	// eventually will support multiple attackers
	for i := 0; i < 1; i++ {
		pms = append(pms, protocol.NewManager(chain))
	}

	return Manager{
		pms:           pms,
		targetDataDir: targetDataDir,
	}
}

func (m *Manager) StartServers() error {
	for _, pm := range m.pms {
		s := p2p.MakeP2PServer(pm)

		err := s.Start()
		if err != nil {
			return err
		}

		p2p.AddLocalPeer(s, m.targetDataDir)

		m.servers = append(m.servers, s)
	}

	return nil
}

func (m *Manager) StopServers() {
	for _, s := range m.servers {
		s.Stop()
	}
}

func (m *Manager) GetRoutinesChannel(idx int) chan attack.Routine {
	return m.pms[idx].Routines
}
