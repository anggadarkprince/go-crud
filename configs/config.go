package configs

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Auth     AuthConfig
	Database DatabaseConfig
	Session  SessionConfig
}

// Global config instance
var Configs *Config

// Load all configs (call this once at app startup)
func Load() (*Config, error) {
	// Set config file name and type
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Allow Viper to read environment variables
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// Don't fail if .env doesn't exist in production
		// Viper will still read from environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	Configs = &Config{
		App:      LoadAppConfig(),
		Auth:     LoadAuthConfig(),
		Database: LoadDatabaseConfig(),
		Session:  LoadSessionConfig(),
	}

	return Configs, nil
}

func Get() *Config {
	return Configs
}
