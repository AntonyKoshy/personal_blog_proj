package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost  string
	DBPort  int
	DBUser  string
	DBPass  string
	DBName  string
	SSLMode string
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // read env vars if config missing

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		DBHost: viper.GetString("DB_HOST"),
		DBPort: viper.GetInt("DB_PORT"),
		DBUser: viper.GetString("DB_USER"),
		DBName: viper.GetString("DB_NAME"),

		// DBPass: viper.GetString("DB_PASSWORD"),
	}, nil
}
