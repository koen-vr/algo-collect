package account

import (
	"os"

	"github.com/algorand/go-algorand-sdk/crypto"
)

// WARN Do Not use in production this is fast and dirty plain text
func saveToFile(path string, acc crypto.Account) error {
	f1, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f1.Close()
	f1.Write(acc.PrivateKey)
	return nil
}

// WARN Do Not use in production this is fast and dirty plain taxt
func loadFromFile(path string) (crypto.Account, error) {
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
