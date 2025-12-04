package configs

import "github.com/spf13/viper"

type AuthConfig struct {
	JwtSecret string
	JwtExpired int
	PersonalToken string
	ResetExpired int
}

func LoadAuthConfig() AuthConfig {
	viper.SetDefault("JWT_SECRET", "jwt-secret")
	viper.SetDefault("JWT_EXPIRED", 7200)
	viper.SetDefault("PERSONAL_TOKEN", "personal-token")
	viper.SetDefault("RESET_EXPIRED", 7200)

	return AuthConfig{
		JwtSecret: viper.GetString("JWT_SECRET"),
		JwtExpired: viper.GetInt("JWT_EXPIRED"),
		PersonalToken: viper.GetString("PERSONAL_TOKEN"),
		ResetExpired: viper.GetInt("RESET_EXPIRED"),
	}
}