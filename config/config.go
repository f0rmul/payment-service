package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
}

func NewConfig(cfgFile string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(cfgFile, cfg); err != nil {
		return nil, errors.Wrap(err, "cleanenv.ReadConfig()")
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, errors.Wrap(err, "cleanenv.UpdateEnv()")
	}
	return cfg, nil
}

// Logger config
type Logger struct {
	Development       bool   `yaml:"development"`
	DisableCaller     bool   `yaml:"disableCaller"`
	DisableStacktrace bool   `yaml:"disableStacktrace:"`
	Encoding          string `yaml:"encoding"`
	Level             string `yaml:"level"`
}

func (l *Logger) IsDevelopment() bool        { return l.Development }
func (l *Logger) IsCallerDisabled() bool     { return l.DisableCaller }
func (l *Logger) IsStackTraceDisabled() bool { return l.DisableStacktrace }
func (l *Logger) GetEncoding() string        { return l.Encoding }
func (l *Logger) GetLevel() string           { return l.Level }
