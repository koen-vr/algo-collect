package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/koen-vr/algo-collect/shared/network"
)

func init() {
	Network.AddCommand(networkStopCmd)
	Network.AddCommand(networkStartCmd)
}

var Network = &cobra.Command{
	Use:   "network",
	Short: "Provides the tools to control a local network node",
	Long:  `A set of commands to support the management of a network nodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
}

var networkStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the network",
	Long:  `Stop the local node and cleans up artifacts.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := network.Destroy(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var networkStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the network",
	Long:  `Start the local node with the provided configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := network.Create(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
