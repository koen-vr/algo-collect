package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	acc "github.com/vecno-io/go-pyteal/account"
)

func init() {
	Account.AddCommand(accountCreateCmd)
}

var Account = &cobra.Command{
	Use:   "account",
	Short: "tools to manage network accounts",
	Long:  `A set of commands to support the management of accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
}

var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates the account",
	Long:  `Creates a secure managment account.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := acc.Create("manager", setup.Passphrase); err != nil {
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
