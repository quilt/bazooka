package attack

import (
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gopkg.in/yaml.v2"
)

type Initialization struct {
	Height   uint64
	Accounts map[common.Address]Account
}

type Account struct {
	Key     hexutil.Bytes
	Balance uint64
	Code    hexutil.Bytes
	Salt    hexutil.Bytes
}

type Attack struct {
	Initialization
	Routines []Routine
}

func (a *Attack) Load(s string) error {
	data, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, a); err != nil {
		return err
	}

	return nil
}

func (a *Attack) Accounts() map[common.Address]Account {
	return a.Initialization.Accounts
}
