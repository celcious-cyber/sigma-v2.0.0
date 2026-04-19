package config

import (
	"os"
)

type Config struct {
	JWTSecret            string
	DBPath               string
	XenditSecretKey      string
	XenditCallbackToken  string
}

var AppConfig *Config

func InitConfig() {
	AppConfig = &Config{
		JWTSecret:           getEnv("JWT_SECRET", "sigma-secret-key-2026"),
		DBPath:              getEnv("DB_PATH", "data/sigma.db"),
		XenditSecretKey:     getEnv("XENDIT_SECRET_KEY", "xnd_development_..."),
		XenditCallbackToken: getEnv("XENDIT_CALLBACK_TOKEN", "sigma-callback-token"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
