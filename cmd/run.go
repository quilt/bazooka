package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lightclient/bazooka/attack"
	"github.com/lightclient/bazooka/simulator"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run an attack against a victim node",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an attack name to execute")
		}

		if attack.IsValidAttack(args[0]) {
			return nil
		}

		return fmt.Errorf("unknown attack specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		setupLogger()

		db := rawdb.NewMemoryDatabase()

		blockchain, err := simulator.InitBlockchain(db)
		if err != nil {
			panic(fmt.Errorf("Error initializing chain: %s", err))
		}

		sm := simulator.NewManager(blockchain, 1)

		runner, err := attack.AttackRunnerFromString(args[0], sm.GetRoutinesChannel(0))
		if err != nil {
			panic(fmt.Errorf("Error initializing attack: %s", err))
		}

		runner.Run()

		sm.StartServers()
		time.Sleep(30 * time.Second)
		sm.StopServers()
	},
}

func setupLogger() {
	var (
		ostream log.Handler
		glogger *log.GlogHandler
	)

	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)

	if usecolor {
		output = colorable.NewColorableStderr()
	}

	ostream = log.StreamHandler(output, log.TerminalFormat(usecolor))
	glogger = log.NewGlogHandler(ostream)
	log.Root().SetHandler(glogger)
	glogger.Verbosity(log.Lvl(5))
}
