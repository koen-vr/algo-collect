package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vecno-io/algo-collection/shared/network"
)

var StopnetCmd = &cobra.Command{
	Use:   "stopnet",
	Short: "Stop the development network",
	Long:  `Stops the development network and cleans out all artifacts.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := network.Destroy(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println()
	},
}
