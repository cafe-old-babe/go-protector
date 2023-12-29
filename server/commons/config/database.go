package config

type database struct {
	Driver          string `yaml:"driver" json:"driver,omitempty"`
	Username        string `yaml:"username" json:"username,omitempty"`
	Password        string `yaml:"password" json:"password,omitempty"`
	Host            string `yaml:"host" json:"host,omitempty"`
	Post            string `yaml:"post" json:"post,omitempty"`
	Dbname          string `yaml:"dbname" json:"dbname,omitempty"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime" json:"connMaxIdleTime,omitempty"`
	ConnMaxLifeTime int    `yaml:"connMaxLifeTime" json:"connMaxLifeTime,omitempty"`
	MaxIdleConns    int    `yaml:"maxIdleConns" json:"maxIdleConns,omitempty"`
	MaxOpenConns    int    `yaml:"maxOpenConns" json:"maxOpenConns,omitempty"`
}
