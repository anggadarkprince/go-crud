package configs

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
    Name string
    Environment string
    Port uint
    Debug bool
}

func LoadAppConfig() AppConfig {
    viper.SetDefault("APP_NAME", "Application")
    viper.SetDefault("APP_ENV", "production")
    viper.SetDefault("APP_PORT", 8080)
    viper.SetDefault("APP_DEBUG", false)
    
    return AppConfig{
        Name: viper.GetString("APP_NAME"),
        Environment: viper.GetString("APP_ENV"),
        Port: viper.GetUint("APP_PORT"),
        Debug: viper.GetBool("APP_DEBUG"),
    }
}