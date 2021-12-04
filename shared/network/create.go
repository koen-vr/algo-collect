package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	exec "github.com/vecno-io/algo-collection/shared/execute"
)

type Config struct {
	Version            uint64
	GossipFanout       uint64
	NetAddress         string
	DNSBootstrapID     string
	EnableProfiler     bool
	EnableDeveloperAPI bool
}

func Create() error {
	fmt.Println("### Creating private network")

	out, err := exec.List([]string{
		"-c", "goal network create -n tn50e -t ./network.json -r ./net1",
	})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return err
	}

	// Update the config to enable the developer api
	// TODO: Fix this hack, the config struct is hacky

	cfg := Config{}
	file, err := ioutil.ReadFile("./net1/Primary/config.json")
	json.Unmarshal(file, &cfg)
	if nil != err {
		return err
	}
	cfg.EnableDeveloperAPI = true

	jsonString, _ := json.Marshal(cfg)
	ioutil.WriteFile("./net1/Primary/config.json", jsonString, os.ModePerm)

	out, err = exec.List([]string{
		"-c", "goal network start -r ./net1",
	})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return err
	}

	// Start the network

	out, err = exec.List([]string{
		"-c", "goal network status -r ./net1",
	})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return err
	}

	return nil
}
