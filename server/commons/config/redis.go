package config

type redis struct {
	Addr     string `yaml:"addr" json:"addr,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
	Username string `yaml:"username" json:"username,omitempty"`
	Db       int    `yaml:"db" json:"db,omitempty"`
}
