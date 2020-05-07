package attack

import (
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lightclient/bazooka/routine"
	"gopkg.in/yaml.v2"
)

type Account struct {
	Key     hexutil.Bytes
	Balance uint64
	Code    hexutil.Bytes
	Salt    hexutil.Bytes
}

type Attack struct {
	Accounts map[common.Address]Account
	Routines []routine.Routine
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
