package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bazooka",
	Short: "Bazooka is an attack orchestration tool targeting Ethereum clients.",
}

func init() {
	runCmd.Flags().StringVarP(&targetDataDir, "target-data-dir", "t", "eth-sim", "geth datadir of target node")
	rootCmd.AddCommand(runCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
