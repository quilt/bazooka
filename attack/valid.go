package attack

import (
	"fmt"
	"os"
	"strings"

	"github.com/lightclient/bazooka/routine"
)

func IsValidAttack(s string) bool {
	_, err := os.Stat(s)
	if err == nil {
		return true
	}

	return true
}

func AttackRunnerFromString(s string, c chan routine.Routine) (*Runner, error) {
	if strings.ToLower(s) == "sample" {
		return NewSampleAttack(c)
	}

	return nil, fmt.Errorf("invalid attack name")
}
