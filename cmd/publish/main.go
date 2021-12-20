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
}

func init() {
	cobra.OnInitialize(onInitialize)

	rootCmd.AddCommand(command.Assets)
	rootCmd.AddCommand(command.Contract)

	rootCmd.AddCommand(command.Account)
	rootCmd.AddCommand(command.Network)
}

func onInitialize() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	viper.SetEnvPrefix("algo")
	viper.BindEnv("type")
	viper.BindEnv("pass")
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
