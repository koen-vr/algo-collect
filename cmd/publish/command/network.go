package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	net "github.com/vecno-io/go-pyteal/network"
)

func init() {
	Network.AddCommand(networkStopCmd)
	Network.AddCommand(networkStartCmd)
	Network.AddCommand(networkStatusCmd)

	Network.AddCommand(networkCreateCmd)
	Network.AddCommand(networkDestroyCmd)
}

var Network = &cobra.Command{
	Use:   "network",
	Short: "tools to control a local network node",
	Long:  `A set of commands to support the management of a network nodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
}

var networkStartCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the network node",
	Long:  `Start the local network node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := net.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

var networkStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stops the network node",
	Long:  `Stop the local network node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := net.Stop(); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

var networkStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "get the network status",
	Long:  `Shows the current network status.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := net.Status(); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

var networkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a network node",
	Long:  `Create and setup a local network node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := net.Create(); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(false); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

var networkDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroys a network node",
	Long:  `Destroy a local network node and clean up.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := net.Destroy(); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}
