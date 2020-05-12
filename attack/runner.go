package attack

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

func (a *Attack) NewRunner(c chan Routine) Runner {
	a.SignAll()
	return Runner{routines: a.Routines, c: c}
}
