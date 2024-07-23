package config

import (
	"github.com/go-yaml/yaml"
	"go-protector/server/internal/consts"
	"os"
	"sync"
)

var _config *config

var lock sync.Mutex

// 2-5	【实战】引入配置文件-掌握GO语言操作文件、Tag特性
type config struct {
	Database database `yaml:"database"`
	Logger   logger   `yaml:"logger"`
	Server   server   `yaml:"server"`
	Redis    redis    `yaml:"redis"`
	Jwt      Jwt      `yaml:"jwt"`
	Email    Email    `yaml:"email"`
}

func GetConfig() *config {
	// 2-6	【实战】玩点儿花活，改写配置文件代码-掌握DCL单例模式
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
