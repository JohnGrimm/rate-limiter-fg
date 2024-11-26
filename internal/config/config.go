package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	RedisHost     string `mapstructure:"redis_host"`
	RedisPort     string `mapstructure:"redis_port"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisDB       int    `mapstructure:"redis_db"`

	DefaultRateLimit   int `mapstructure:"default_rate_limit"`
	DefaultExpiry      int `mapstructure:"default_expiry"`
	DefaultTimeBlocked int `mapstructure:"default_time_blocked"`
}

var cfg Config

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml") // ou "json", "toml", etc.

	viper.AutomaticEnv() // lê variáveis de ambiente

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	log.Printf("load config, %v", cfg)
}

func GetConfig() Config {
	return cfg
}
