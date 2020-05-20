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
			return errors.New("requires an attack yaml to execute")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		defer func() {
			if err != nil {
				log.Error(fmt.Sprintf("%s", err))
			}
		}()

		setupLogger()

		var attack attack.Attack

		err = attack.Load(args[0])
		if err != nil {
			return
		}

		blockchain, err := simulator.InitBlockchain(rawdb.NewMemoryDatabase(), attack.Accounts)
		if err != nil {
			return
		}

		sm := simulator.NewManager(blockchain)

		runner, err := attack.NewRunner(blockchain, sm.GetRoutinesChannel(0))
		if err != nil {
			return
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
