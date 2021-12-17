package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	acc "github.com/vecno-io/go-pyteal/account"

	"github.com/algorand/go-algorand-sdk/mnemonic"
)

var accountName string

func init() {
	Account.AddCommand(accountInfoCmd)
	Account.AddCommand(accountCreateCmd)

	Account.PersistentFlags().StringVarP(&accountName, "name", "n", "account", "name of license for the project (default: account)")
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

var accountInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get account info",
	Long:  `Loads an account from file and shows the info when it exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := acc.Load(accountName, setup.Passphrase)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		recovery, err := mnemonic.FromPrivateKey(acc.PrivateKey)
		if nil != err {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		fmt.Printf("> Info: %s.acc loaded. Make sure to backup the file or safe the info below. \n", accountName)
		fmt.Println(">> Account: ", acc.Address.String())
		fmt.Println(">> Recovery: ", recovery)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates the account",
	Long:  `Creates a secure managment account.`,
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := acc.Create(accountName, setup.Passphrase)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		recovery, err := mnemonic.FromPrivateKey(acc.PrivateKey)
		if nil != err {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		fmt.Printf("> Info: %s.acc created. Make sure to backup the file or safe the info below. \n", accountName)
		fmt.Println(">> Account: ", acc.Address.String())
		fmt.Println(">> Recovery: ", recovery)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}
