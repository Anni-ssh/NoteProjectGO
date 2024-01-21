package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type StartupConfig struct {
	Env          string `yaml:"env"`
	DataBasePath string `yaml:"dataBasePath"`
}

func CreateCfg(path string) (*StartupConfig, error) {
	const operation = "config.CreateCfg"
	var cfg StartupConfig

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", operation, err)
	}
	return &cfg, nil
}
