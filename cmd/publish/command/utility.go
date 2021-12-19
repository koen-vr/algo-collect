package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	cfg "github.com/vecno-io/go-pyteal/config"
)

var setup = cfg.Setup{}

func onInitialize(validate bool) error {
	if err := viper.Unmarshal(&setup); nil != err {
		return err
	}
	if err := cfg.OnInitialize(setup, validate); nil != err {
		return err
	}
	return nil
}

func getListOfFiles(ext, path string) ([]string, error) {
	list := make([]string, 0)
	info, err := os.Stat(path)
	if nil != err {
		return list[:], fmt.Errorf("get file list: %s", err)
	}
	if !info.IsDir() {
		return list[:], fmt.Errorf("get file list:: not a direcotry: %s", path)
	}

	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ext {
			list = append(list, path)
		}
		return nil
	}); err != nil {
		return list[:], fmt.Errorf("get file list: walk tree: %s", err)
	}

	return list[:], nil
}

func getUnitName(nr uint32) string {
	str := fmt.Sprintf("%%s#%%0%dd", len(viper.GetString("META_COLLECT_MAXCOUNT")))
	return fmt.Sprintf(str, viper.GetString("META_COLLECT_TAG"), nr)
}

func getApplicationId(name string) (string, error) {
	id, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.id", viper.GetString("DATA"), name))
	if err != nil {
		return "", err
	}
	return strings.Split(string(id), "\n")[0], nil
}

func getApplicationUrl(name string) (string, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json.pin", viper.GetString("DATA"), name))
	if err != nil {
		return "", err
	}
	pin := PinFileResponse{}
	if err := json.Unmarshal(data, &pin); nil != err {
		return "", err
	}
	return fmt.Sprintf("ipfs://%s", pin.IpfsHash), nil
}
