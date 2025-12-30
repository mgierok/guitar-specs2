package config

import "github.com/spf13/viper"

type Config struct {
	Env        string
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() (Config, error) {
	v := viper.New()
	v.SetDefault("ENV", "development")
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "guitar_specs")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetEnvPrefix("GUITAR_SPECS")
	v.AutomaticEnv()

	cfg := Config{
		Env:        v.GetString("ENV"),
		ServerPort: v.GetString("SERVER_PORT"),
		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetString("DB_PORT"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBName:     v.GetString("DB_NAME"),
		DBSSLMode:  v.GetString("DB_SSLMODE"),
	}

	return cfg, nil
}
