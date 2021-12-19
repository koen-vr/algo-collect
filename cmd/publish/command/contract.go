package command

import (
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	acc "github.com/vecno-io/go-pyteal/account"
	app "github.com/vecno-io/go-pyteal/contract"
)

func init() {
	Contract.AddCommand(contractBuildCmd)
	Contract.AddCommand(contractDeployCmd)
}

var Contract = &cobra.Command{
	Use:   "contract",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
}

var contractBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.Build([]string{
			viper.GetString("APP_PROG"),
			viper.GetString("APP_CLEAR"),
		}); nil != err {
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

var contractDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := acc.Load(viper.GetString("APP_MANAGER"), setup.Passphrase)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		// Verify balance here?
		// There is no verification
		// and we manualy set schema.
		// So there is no balance info.

		if err := app.Deploy(app.Setup{
			Manager:      acc,
			ClearProg:    viper.GetString("APP_CLEAR"),
			ApprovalProg: viper.GetString("APP_PROG"),
			LocalSchema: types.StateSchema{
				NumUint:      1,
				NumByteSlice: 0,
			},
			GlobalSchema: types.StateSchema{
				NumUint:      0,
				NumByteSlice: 64,
			},
		}); nil != err {
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
