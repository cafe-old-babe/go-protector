package config

type Jwt struct {
	TokenPre       string `yaml:"tokenPre"`
	TokenTimeout   int    `yaml:"tokenTimeout"`
	SessionTimeout int    `yaml:"sessionTimeout"`
}
