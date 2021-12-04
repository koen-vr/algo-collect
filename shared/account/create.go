package account

import (
	"errors"
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/crypto"

	exec "github.com/vecno-io/algo-collection/shared/execute"
)

func Create() crypto.Account {
	acc := crypto.GenerateAccount()
	return acc
}

func CreateAndFund(algo uint, from string) (crypto.Account, error) {
	acc := Create()
	if algo > 10 {
		algo = 10
	}
	if !IsAlgorandAddress(from) {
		return acc, errors.New("failed to fund from invalid address")
	}
	fmt.Println("\n### Funding new account")
	out, err := exec.List([]string{"-c", fmt.Sprintf(
		"goal -d ./net1/Primary clerk send -a %d -f %s -t %s",
		algo*10000000000000, from, acc.Address.String(),
	)})
	fmt.Println(out)
	if nil != err {
		fmt.Println(err)
	}
	return acc, nil
}

// WARN Do Not use in production this is fast and dirty plain taxt
func SaveToFile(path string, acc crypto.Account) error {
	f1, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f1.Close()
	f1.Write(acc.PrivateKey)
	return nil
}

// WARN Do Not use in production this is fast and dirty plain taxt
func LoadFromFile(path string) (crypto.Account, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return crypto.Account{}, err
	}

	ac, err := crypto.AccountFromPrivateKey(data)
	if err != nil {
		return crypto.Account{}, err
	}
	return ac, nil
}
