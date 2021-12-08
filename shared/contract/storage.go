package contract

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func saveToFile(path string, id uint64) error {
	str, err := json.Marshal(id)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, str, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func loadFromFile(path string) (uint64, error) {
	id := uint64(0)
	dat, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(dat[:], &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
