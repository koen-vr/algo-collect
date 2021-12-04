package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vecno-io/algo-collection/cmd/collection/cmd"
)

var (
	cfgFile     string
	userLicense string
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "Run development and test utilities",
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cmd.StopnetCmd)
	rootCmd.AddCommand(cmd.StartnetCmd)

	rootCmd.AddCommand(cmd.CreateappCmd)
	rootCmd.AddCommand(cmd.CreateAssetCmd)
	rootCmd.AddCommand(cmd.CreateWalletCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".algo-collection")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// func main() {
// 	var err error
// 	defer func() {
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 	}()

// 	fmt.Println("### Creating private network")
// 	network.Create()
// 	defer func() {
// 		fmt.Println("### Destroying private network")
// 		network.Destroy()
// 	}()

// 	fmt.Println("\n### Compile approval prog")
// 	out, err := exec.List([]string{
// 		"-c", "python3 ./contracts/collection.py > ./contracts/collection.teal",
// 	})
// 	fmt.Println(out)
// 	if nil != err {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("\n### Compile clear prog")
// 	out, err = exec.List([]string{
// 		"-c", "python3 ./contracts/clear.py > ./contracts/clear.teal",
// 	})
// 	fmt.Println(out)
// 	if nil != err {
// 		fmt.Println(err)
// 	}
// 	// cmd := fmt.Sprintf(
// 	// 	"goal app create --creator %s --approval-prog %s",
// 	// 	acc.Address.String()
// 	// )
// 	acc := account.FirstGoalAccount()
// 	cmdList := []string{
// 		"goal app create",
// 		fmt.Sprintf(" --creator %s", acc),
// 		fmt.Sprintf(" --approval-prog %s", "./contracts/collection.teal"),
// 		fmt.Sprintf(" --clear-prog %s", "./contracts/clear.teal"),
// 		fmt.Sprintf(" --global-byteslices %d", 64),
// 		fmt.Sprintf(" --global-ints %d", 0),
// 		fmt.Sprintf(" --local-byteslices %d", 0),
// 		fmt.Sprintf(" --local-ints %d", 1),
// 		" --on-completion OptIn",
// 		" -d ./net1/Primary |",
// 		" grep Created | awk '{ print $6 }'",
// 	}
// 	var b strings.Builder
// 	for _, p := range cmdList {
// 		fmt.Fprintf(&b, "%s", p)
// 	}
// 	fmt.Println("\n### Creating Application")
// 	fmt.Println(b.String())
// 	out, err = exec.List([]string{
// 		"-c", b.String(),
// 	})
// 	fmt.Println(out)
// 	if nil != err {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("\n### Check Application")

// 	out, err = exec.List([]string{
// 		"-c", fmt.Sprintf("goal app read -d ./net1/Primary --app-id 1 --guess-format --local --from %s", acc),
// 	})
// 	fmt.Println(out)
// 	if nil != err {
// 		fmt.Println(err)
// 	}
// 	out, err = exec.List([]string{
// 		"-c", "goal app read -d ./net1/Primary --app-id 1 --guess-format --global",
// 	})
// 	fmt.Println(out)
// 	if nil != err {
// 		fmt.Println(err)
// 	}

// 	fmt.Print("\nPress 'Enter' to continue...")
// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
// }
