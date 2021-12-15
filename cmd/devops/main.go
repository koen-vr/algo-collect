package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/koen-vr/algo-collect/cmd/devops/cmd"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publication tool for algo-collet",
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cmd.Create)
	rootCmd.AddCommand(cmd.Network)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.algo-cfg)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".algo-cfg"
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".algo-cfg")
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
