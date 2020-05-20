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

func (a *Attack) NewRunner(bc *core.BlockChain, c chan Routine) (*Runner, error) {
	err := a.SignAndAssemble(bc)
	if err != nil {
		return nil, err
	}

	return &Runner{routines: a.Routines, c: c}, nil
}
