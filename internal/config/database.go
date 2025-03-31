package config

import "fmt"

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode" default:"disable"`
}

func (c *Database) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
	// return fmt.Sprintf(
	// 	"host=%s port=%d user=%s password=%s dbname=%s sslmode=%v",
	// 	c.Host,
	// 	c.Port,
	// 	c.User,
	// 	c.Password,
	// 	c.DBName,
	// 	c.SSLMode,
	// )
}
