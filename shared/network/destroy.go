package network

import (
	"fmt"

	exec "github.com/vecno-io/algo-collection/shared/execute"
)

func Destroy() error {
	fmt.Println("### Destroying private network")

	out, err := exec.List([]string{
		"-c", "goal network stop -r ./net1",
	})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return err
	}

	out, err = exec.List([]string{
		"-c", "goal network delete -r ./net1",
	})
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
