package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	net "github.com/vecno-io/go-pyteal/network"
)

func init() {
	// Network.AddCommand(networkStopCmd)
	// Network.AddCommand(networkStartCmd)
	Network.AddCommand(networkCreateCmd)
	Network.AddCommand(networkDestroyCmd)
}

func getNetConfig() net.Config {
	path, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return net.Config{
		Type:     net.Devnet,
		NodePath: fmt.Sprintf("%s/node", path),
		DataPath: fmt.Sprintf("%s/node/devnet-data", path),
	}
	// return net.Config{
	// 	Type:     net.Testnet,
	// 	NodePath: fmt.Sprintf("%s/node", path),
	// 	DataPath: fmt.Sprintf("%s/node/testnet-data", path),
	// }
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

// var networkStartCmd = &cobra.Command{
// 	Use:   "start",
// 	Short: "Starts the network node",
// 	Long:  `Start the local network node.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if err := net.Start(getNetConfig()); err != nil {
// 			fmt.Fprintln(os.Stderr, "error: "+err.Error())
// 			os.Exit(1)
// 		}
// 	},
// }

// var networkStopCmd = &cobra.Command{
// 	Use:   "stop",
// 	Short: "Stops the network node",
// 	Long:  `Stop the local network node.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if err := net.Stop(getNetConfig()); err != nil {
// 			fmt.Fprintln(os.Stderr, "error: "+err.Error())
// 			os.Exit(1)
// 		}
// 	},
// }

var networkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a network node",
	Long:  `Create and setup a local network node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := net.Create(getNetConfig()); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

var networkDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys a network node",
	Long:  `Destroy a local network node and clean up.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := net.Destroy(getNetConfig()); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}
