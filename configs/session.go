package configs

import "github.com/spf13/viper"

type SessionConfig struct {
	StoreName  string
	CookieName string
	Lifetime int
	Secret string
	Path string
	Domain string
	Secure bool
	SameSite string
}

func LoadSessionConfig() SessionConfig {
	viper.SetDefault("SESSION_STORE_NAME", "session_store")
	viper.SetDefault("SESSION_COOKIE", "session")
	viper.SetDefault("COOKIE_SECRET", "secret")
	viper.SetDefault("COOKIE_LIFETIME", 7200)
	viper.SetDefault("COOKIE_PATH", "/")
	viper.SetDefault("COOKIE_DOMAIN", "localhost")
	viper.SetDefault("COOKIE_SECURE", false)
	viper.SetDefault("COOKIE_SAME_SITE", "lax")

	return SessionConfig{
		StoreName: viper.GetString("SESSION_STORE_NAME"),
		CookieName: viper.GetString("SESSION_COOKIE"),
		Secret: viper.GetString("COOKIE_SECRET"),
		Lifetime: viper.GetInt("COOKIE_LIFETIME"),
		Path: viper.GetString("COOKIE_PATH"),
		Domain: viper.GetString("COOKIE_DOMAIN"),
		Secure: viper.GetBool("COOKIE_SECURE"),
		SameSite: viper.GetString("COOKIE_SAME_SITE"),
	}
}