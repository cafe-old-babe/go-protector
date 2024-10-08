package config

import (
	"github.com/go-yaml/yaml"
	"go-protector/server/internal/consts"
	"os"
	"sync"
)

var _config *config

var lock sync.Mutex

type config struct {
	Database database `yaml:"database"`
	Logger   logger   `yaml:"logger"`
	Server   server   `yaml:"server"`
	Redis    redis    `yaml:"redis"`
	Jwt      Jwt      `yaml:"jwt"`
	Email    Email    `yaml:"email"`
}

func GetConfig() *config {
	if _config != nil {
		return _config
	}
	lock.Lock()
	defer lock.Unlock()
	if _config != nil {
		return _config
	}

	_config = &config{}

	configPath := os.Getenv(consts.EnvConfig)
	data, _ := os.ReadFile(configPath)

	if err := yaml.Unmarshal(data, _config); err != nil {
		panic(err)
	}
	return _config
}
