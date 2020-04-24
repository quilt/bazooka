package simulator

import (
	"github.com/ethereum/go-ethereum/core"
	ep2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/lightclient/bazooka/attack"
	"github.com/lightclient/bazooka/p2p"
	"github.com/lightclient/bazooka/protocol"
	"github.com/lightclient/bazooka/routine"
)

type Manager struct {
	pms     []*protocol.Manager
	servers []*ep2p.Server
	attack  attack.Runner
}

func NewManager(chain *core.BlockChain, n int) Manager {
	var pms []*protocol.Manager
	for i := 0; i < n; i++ {
		pms = append(pms, protocol.NewManager(chain))
	}

	return Manager{
		pms: pms,
	}
}

func (m *Manager) StartServers() error {
	for _, pm := range m.pms {
		s := p2p.MakeP2PServer(pm)

		err := s.Start()
		if err != nil {
			return err
		}

		p2p.AddLocalPeer(s)

		m.servers = append(m.servers, s)
	}

	return nil
}

func (m *Manager) StopServers() {
	for _, s := range m.servers {
		s.Stop()
	}
}

func (m *Manager) GetRoutinesChannel(idx int) chan routine.Routine {
	return m.pms[idx].Routines
}
