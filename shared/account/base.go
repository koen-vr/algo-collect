package account

import (
	"github.com/algorand/go-algorand-sdk/types"

	exec "github.com/vecno-io/algo-collection/shared/execute"
)

func FirstGoalAccount() (string, error) {
	out, err := exec.List([]string{"-c", "goal account list -d ./net1/Primary | awk '{ print $3 }' | head -n 1"})
	if nil != err {
		return "", err
	}
	return string(out), nil
}

func SecondGoalAccount() (string, error) {
	out, err := exec.List([]string{"-c", "goal account list -d ./net1/Primary | awk '{ print $3 }' | tail -1"})
	if nil != err {
		return "", err
	}
	return string(out), nil
}

func IsAlgorandAddress(addr string) bool {
	if _, err := types.DecodeAddress(addr); nil != err {
		return false
	}
	return true
}
