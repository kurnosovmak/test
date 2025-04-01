package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"app"`
	Database Database `mapstructure:"db"`
	JWT      JWT      `mapstructure:"jwt"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
