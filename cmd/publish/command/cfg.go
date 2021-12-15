package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cfg "github.com/vecno-io/go-pyteal/config"
)

var (
	NetTarget string
)

func getNetConfig() cfg.Network {
	path, err := os.UserHomeDir()
	cobra.CheckErr(err)

	c := cfg.Network{}
	c.NodePath = fmt.Sprintf("%s/node", path)

	switch NetTarget {
	case "devnet":
		c.Type = cfg.Devnet
		c.DataPath = fmt.Sprintf("%s/node/devnet-data", path)
	case "testnet":
		c.Type = cfg.Testnet
		c.DataPath = fmt.Sprintf("%s/node/testnet-data", path)
	case "mainnet":
		c.Type = cfg.Mainnet
		c.DataPath = fmt.Sprintf("%s/node/mainnet-data", path)
	default:
		err := fmt.Errorf("get config: unknown target: %s", NetTarget)
		fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}

	return c
}
