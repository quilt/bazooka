package attack

import (
	"fmt"
	"strings"

	"github.com/lightclient/bazooka/routine"
)

var ATTACKS = []string{"sample"}

func IsValidAttack(s string) bool {
	for _, a := range ATTACKS {
		if a == s {
			return true
		}
	}

	return false
}

func AttackRunnerFromString(s string, c chan routine.Routine) (*Runner, error) {
	if strings.ToLower(s) == "sample" {
		return NewSampleAttack(c)
	}

	return nil, fmt.Errorf("invalid attack name")
}
