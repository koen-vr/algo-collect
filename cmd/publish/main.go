package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/koen-vr/algo-collect/cmd/publish/command"
)

var (
	cfgFile string
)

type Config struct {
	DataPath string `mapstructure:"ALGORAND_DATA"`
}

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "Run development and test utilities",
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
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(command.Network)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.env)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Look for .env file
		viper.AddConfigPath(".")
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
	}

	viper.AutomaticEnv()
	// Override loaded environment variables
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
