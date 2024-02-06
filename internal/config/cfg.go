package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type StartupConfig struct {
	Env          string `json:"env"`
	DataBasePath string `json:"dataBasePath"`
	PagesPath    string `json:"pagesPath"`
}

func CreateCfg(configPath string) (*StartupConfig, error) {
	const operation = "config.CreateCfg"
	var cfg StartupConfig

	//Проверка корректности пути
	if configPath == "" {
		return nil, fmt.Errorf("%s: Пустой путь к конфигу", operation)
	}

	//	проверка наличия файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: Не найден конфигурационный файл: %s", operation, configPath)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", operation, err)
	}

	return &cfg, nil
}
