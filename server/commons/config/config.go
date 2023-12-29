package config

var Config = &config{}

type config struct {
	Database database `yaml:"database"`
	Logger   logger   `yaml:"logger"`
	Server   server   `yaml:"server"`
	Redis    redis    `yaml:"redis"`
	Jwt      Jwt      `yaml:"jwt"`
}
