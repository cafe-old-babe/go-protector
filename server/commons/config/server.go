package config

type server struct {
	Port       string `yaml:"port" json:"port,omitempty"`
	Env        string `yaml:"env" json:"env,omitempty"`
	RecordPath string `yaml:"recordPath" json:"recordPath,omitempty"`
}
