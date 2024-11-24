package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string
	RedisAddress  string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func init() {
	viper.SetDefault("ServerAddress", ":8080")
	viper.SetDefault("RedisAddress", "redis:6379")
	viper.SetDefault("RedisPassword", "")
	viper.SetDefault("RedisDB", 0)
}
