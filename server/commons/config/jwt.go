package config

type Jwt struct {
	TokenPre       string `yaml:"tokenPre"`
	Timeout        int    `yaml:"timeout"`
	SessionTimeout int    `yaml:"sessionTimeout"`
}
