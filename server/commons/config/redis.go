package config

type redis struct {
	Addr     string `yaml:"addr" json:"addr"`
	Password string `yaml:"password" json:"password"`
	Username string `yaml:"username" json:"username"`
	Db       int    `yaml:"db" json:"db"`
}
