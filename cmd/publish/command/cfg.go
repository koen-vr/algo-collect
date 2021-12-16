package command

import (
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
