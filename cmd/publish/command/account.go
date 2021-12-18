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
	Account.AddCommand(accountKeyCmd)
	Account.AddCommand(accountInfoCmd)
	Account.AddCommand(accountCreateCmd)

	Account.PersistentFlags().StringVarP(&accountName, "name", "n", "manager", "name of license for the project (default: account)")
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

var accountKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "shows account recovery key",
	Long:  `Loads an account from file and shows the recovery key if it exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, err := acc.Load(accountName, setup.Passphrase)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		recovery, err := mnemonic.FromPrivateKey(ac.PrivateKey)
		if nil != err {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		fmt.Printf("> Info: %s.acc loaded. Make sure to backup the file or safe the info below. \n", accountName)
		fmt.Println(">> Account: ", ac.Address.String())
		fmt.Println(">> Recovery: ", recovery)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

var accountInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "show account info",
	Long:  `Shows the account info and balance (and on devnet it tries to fund it).`,
	Run: func(cmd *cobra.Command, args []string) {
		info, err := acc.Info(accountName, setup.Passphrase)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		fmt.Printf("> Info: %s.acc balance. \n", accountName)
		if info.Amount == 0 {
			fmt.Println(">> Make sure to add enough funds for app and asset creation.")
		}

		fmt.Printf(">> Address: %s\n", info.Address)
		fmt.Printf(">> Ballance: %d\n", info.Amount)
		fmt.Printf(">> Assets: %d\n", len(info.Assets))
		fmt.Printf(">> Created Apps: %d\n", len(info.CreatedApps))
		fmt.Printf(">> Created Assets: %d\n", len(info.CreatedAssets))

		if info.Amount == 0 && setup.Target == "devnet" {
			if err := acc.DevFunding(info.Address, 1000000000000); err != nil {
				fmt.Fprintln(os.Stderr, "error: "+err.Error())
				os.Exit(1)
			}
		}
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
		ac, err := acc.Create(accountName, setup.Passphrase)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		recovery, err := mnemonic.FromPrivateKey(ac.PrivateKey)
		if nil != err {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		fmt.Printf("> Info: %s.acc created. Make sure to backup the file or safe the info below. \n", accountName)
		fmt.Println(">> Account: ", ac.Address.String())
		fmt.Println(">> Recovery: ", recovery)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}
