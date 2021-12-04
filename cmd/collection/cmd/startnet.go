package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vecno-io/algo-collection/shared/network"
)

var StartnetCmd = &cobra.Command{
	Use:   "startnet",
	Short: "Start the development network",
	Long:  `start the development network using the network.json config.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := network.Create(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
