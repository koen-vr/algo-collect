package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vecno-io/algo-collection/shared/account"
)

var CreateWalletCmd = &cobra.Command{
	Use:   "createwallet",
	Short: "Create the and backup a wallet",
	Long:  `Creates a wallet for use with other commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("### Creating the collection application")

		fmt.Println("--- Load network accounts")
		fa, err := account.FirstGoalAccount()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		sa, err := account.SecondGoalAccount()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Printf("First network account: %s\n", fa)
		fmt.Printf("Second network account: %s\n", sa)
		fmt.Println()

		fmt.Println("--- Create and fund local accounts")
		ac1, err := account.CreateAndFund(25, fa)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := account.SaveToFile("./ac1.frag", ac1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("First local account: ", ac1.Address)

		ac2, err := account.CreateAndFund(25, sa)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := account.SaveToFile("./ac2.frag", ac2); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Second local account: ", ac2.Address)
	},
}
