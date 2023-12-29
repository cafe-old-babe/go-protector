package initialize

import (
	"github.com/go-yaml/yaml"
	"go-protector/server/commons/config"
	"os"
)

func initConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config.Config)
	return err
}
