package config

type server struct {
	Port       int    `yaml:"port" json:"port,omitempty"`
	Model      string `yaml:"model" json:"model,omitempty"` // debug,release,test
	Env        string `yaml:"env" json:"env,omitempty"`
	RecordPath string `yaml:"recordPath" json:"recordPath,omitempty"`
}
