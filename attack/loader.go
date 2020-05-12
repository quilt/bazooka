package attack

import (
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func (a *Attack) SignAll() {
	for i, r := range a.Routines {
		if r.Ty == SendTxs {
			txs := make([]*types.Transaction, 0)

			for _, tx := range r.Transactions {
				priv, err := crypto.ToECDSA(a.Accounts[tx.From].Key)
				unsignedTx := tx.toEthType()

				signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(big.NewInt(1337)), priv)
				if err != nil {
					panic(err)
				}

				txs = append(txs, signedTx)
			}

			a.Routines[i].SignedTransactions = txs
			fmt.Printf("we got %d unsigned txs and %d signed txs\n", len(r.Transactions), len(a.Routines[i].SignedTransactions))
		}
	}
}
