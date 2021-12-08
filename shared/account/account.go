package account

import (
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"

	exec "github.com/vecno-io/algo-collection/shared/execute"
	net "github.com/vecno-io/algo-collection/shared/network"
)

func Setup() error {
	fmt.Println("### Creating the collection application")

	fmt.Println("--- Load network account")
	sed, err := getSeedAccount()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("- Account sed: %s\n", sed)

	fmt.Println("--- Create local accounts")
	ac1 := crypto.GenerateAccount()
	ac2 := crypto.GenerateAccount()

	if err := saveToFile(fmt.Sprintf("%s/acc1.frag", net.Path()), ac1); err != nil {
		return fmt.Errorf("account: failed to save account nr1: %s", err)
	}
	if err := saveToFile(fmt.Sprintf("%s/acc2.frag", net.Path()), ac2); err != nil {
		return fmt.Errorf("account: failed to save account nr2: %s", err)
	}

	fmt.Printf("- Account nr1: %s\n", ac1.Address.String())
	fmt.Printf("- Account nr2: %s\n", ac2.Address.String())

	fmt.Println("--- Fund local accounts")
	if err := sendFromSeedAccount(10, sed, ac1.Address.String()); err != nil {
		return fmt.Errorf("account: failed to fund account nr1: %s", err)
	}
	if err := sendFromSeedAccount(10, sed, ac2.Address.String()); err != nil {
		return fmt.Errorf("account: failed to fund account nr2: %s", err)
	}

	return nil
}

func IsAlgorandAddress(addr string) bool {
	if _, err := types.DecodeAddress(addr); nil != err {
		return false
	}
	return true
}

func GetMainAccount() (crypto.Account, error) {
	return loadFromFile(fmt.Sprintf("%s/acc1.frag", net.Path()))
}

func GetUserAccount() (crypto.Account, error) {
	return loadFromFile(fmt.Sprintf("%s/acc2.frag", net.Path()))
}

func getSeedAccount() (string, error) {
	out, err := exec.List([]string{"-c", fmt.Sprintf(
		"goal account list -d %s | awk '{ print $3 }' | head -n 1",
		net.NodePath(),
	)})
	if nil != err {
		return "", err
	}
	return string(out), nil
}

func sendFromSeedAccount(amount uint64, from, to string) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %d", amount)
	}
	if !IsAlgorandAddress(to) {
		return fmt.Errorf("invalid to address: %s", to)
	}
	if !IsAlgorandAddress(from) {
		return fmt.Errorf("invalid from address: %s", from)
	}

	out, err := exec.List([]string{"-c", fmt.Sprintf(
		"goal -d %s clerk send -a %d -f %s -t %s",
		net.NodePath(), amount*10000000000000, from, to,
	)})

	fmt.Println(out)
	return err
}
