package config

import "time"

type JWT struct {
	SecretKey     string        `mapstructure:"secret"`
	TokenDuration time.Duration `mapstructure:"expiration"`
}
