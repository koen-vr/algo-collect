package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/koen-vr/algo-collect/shared/account"
	"github.com/koen-vr/algo-collect/shared/asset"
	"github.com/koen-vr/algo-collect/shared/contract"
)

func init() {
	Create.AddCommand(createAppCmd)
	Create.AddCommand(createAssetCmd)
	Create.AddCommand(createAccountCmd)
}

var Create = &cobra.Command{
	Use:   "create",
	Short: "Provides the tools to create network objects",
	Long:  `A set of commands to support the creation of network objects.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
}

var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Deploy application contracts",
	Long:  `Compiles the applications contracts and deployes them to the network.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := contract.Deploy(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var createAssetCmd = &cobra.Command{
	Use:   "asset",
	Short: "Create a new asset in the application",
	Long:  `Configures a new asset and register it within the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := asset.Setup(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var createAccountCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Create new accounts for the application",
	Long:  `Creates new accounts and funds them when funding is available.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := account.Setup(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
