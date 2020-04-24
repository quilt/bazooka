package attack

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/lightclient/bazooka/routine"
)

func NewSampleAttack(c chan routine.Routine) (*Runner, error) {
	var routines []routine.Routine

	coinbaseKey, err := crypto.HexToECDSA("ad0f3019b6b8634c080b574f3d8a47ef975f0e4b9f63e82893e9a7bb59c2d609")
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(0, common.BigToAddress(big.NewInt(0)), big.NewInt(100), params.TxGas, big.NewInt(10), nil)
	signer := types.NewEIP155Signer(big.NewInt(1337))
	signedTx, err1 := types.SignTx(tx, signer, coinbaseKey)
	if err1 != nil {
		panic(err1)
	}

	routines = append(routines, routine.Routine{
		Ty:           routine.SendTxs,
		Transactions: []types.Transaction{*signedTx},
	})

	runner := Runner{
		routines: routines, c: c,
	}

	return &runner, nil
}
