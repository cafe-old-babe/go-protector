package config

type server struct {
	Port       int    `yaml:"port" json:"port,omitempty"`
	Model      string `yaml:"model" json:"model,omitempty"` // debug,release,test
	Env        string `yaml:"env" json:"env,omitempty"`
	Sm4Key     string `yaml:"sm4Key" json:"sm4Key,omitempty"`
	RecordPath string `yaml:"recordPath" json:"recordPath,omitempty"`
	TempPath   string `yaml:"tempPath" json:"tempPath,omitempty"`
}
