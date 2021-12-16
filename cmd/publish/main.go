package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/koen-vr/algo-collect/cmd/publish/command"
)

var rootCmd = &cobra.Command{
	Use:   "publisher",
	Short: "A utility to aid with publishing ARC3 Collections to the Algorand network.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("./contracts"); err != nil {
			fmt.Fprintln(os.Stderr, "contracts folder not available")
			os.Exit(1)
		}
		if _, err := os.Stat("./network.json"); err != nil {
			fmt.Fprintln(os.Stderr, "network setup file not available")
			os.Exit(1)
		}
	},
}

func init() {
	cobra.OnInitialize(onInitialize)

	rootCmd.AddCommand(command.Account)
	rootCmd.AddCommand(command.Network)
}

func onInitialize() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	viper.SetEnvPrefix("algorand")
	viper.BindEnv("node")
	viper.BindEnv("data")

	viper.ReadInConfig()
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
