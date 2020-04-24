package attack

import (
	"github.com/lightclient/bazooka/routine"
)

type Runner struct {
	routines []routine.Routine
	c        chan routine.Routine
}

func (r *Runner) Run() error {
	for _, routine := range r.routines {
		r.c <- routine
	}

	return nil
}
