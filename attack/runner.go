package attack

import "github.com/ethereum/go-ethereum/core"

type Runner struct {
	routines []Routine
	c        chan Routine
}

func (r *Runner) Run() error {
	for _, routine := range r.routines {
		r.c <- routine
	}

	return nil
}

func (a *Attack) NewRunner(bc *core.BlockChain, c chan Routine) Runner {
	a.SignAll()
	a.MakeBlocks(bc)

	return Runner{routines: a.Routines, c: c}
}
