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

	out, err := exec.List([]string{"-c", fmt.Sprintf(
		"goal network create -n tn50e -t ./network.json -r %s", netPath,
	)})
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
	file, err := ioutil.ReadFile(fmt.Sprintf(
		"%s/config.json", nodePath,
	))
	json.Unmarshal(file, &cfg)
	if nil != err {
		return err
	}
	cfg.EnableDeveloperAPI = true

	jsonString, _ := json.Marshal(cfg)

	if ioutil.WriteFile(fmt.Sprintf(
		"%s/config.json", nodePath,
	), jsonString, os.ModePerm); nil != err {
		return err
	}
	// Start the network

	out, err = exec.List([]string{"-c", fmt.Sprintf(
		"goal network start -r %s", netPath,
	)})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return err
	}

	out, err = exec.List([]string{"-c", fmt.Sprintf(
		"goal network status -r %s", netPath,
	)})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return err
	}

	return nil
}

func Destroy() error {
	fmt.Println("### Destroying private network")

	out, err := exec.List([]string{"-c", fmt.Sprintf(
		"goal network stop -r %s", netPath,
	)})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return err
	}

	out, err = exec.List([]string{"-c", fmt.Sprintf(
		"goal network delete -r %s", netPath,
	)})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return err
	}

	exec.List([]string{"-c", "rm -f ./*.rej"})
	exec.List([]string{"-c", "rm -f ./*.txn"})
	exec.List([]string{"-c", "rm -f ./*.txs"})
	exec.List([]string{"-c", "rm -f ./*.frag"})

	return nil
}
