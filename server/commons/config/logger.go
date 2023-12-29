package config

import "gopkg.in/natefinch/lumberjack.v2"

type logger struct {
	Path     string `yaml:"path"  json:"path,omitempty"`
	FileName string `yaml:"fileName" json:"fileName,omitempty"`
	Level    string `yaml:"level" json:"level,omitempty"`
	lumberjack.Logger
}
