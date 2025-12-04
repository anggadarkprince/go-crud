package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     uint
	User     string
	Password string
	Database string
}

func LoadDatabaseConfig() DatabaseConfig {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_USERNAME", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_DATABASE", "sandbox")

	return DatabaseConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetUint("DB_PORT"),
		User:     viper.GetString("DB_USERNAME"),
		Password: viper.GetString("DB_PASSWORD"),
		Database: viper.GetString("DB_DATABASE"),
	}
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Database,
	)
}
