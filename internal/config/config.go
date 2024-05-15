package config

import "github.com/spf13/viper"

func InitConfig(path, name string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	return viper.ReadInConfig()
}
