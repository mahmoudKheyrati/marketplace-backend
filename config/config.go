package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	Postgres PostgresConfig
	Port     int
}
type PostgresConfig struct {
	Host          string
	Port          int
	Username      string
	Password      string
	Database      string
	MaxConnection int
}

var once sync.Once
var config *Config

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	once.Do(func() {
		viper.SetConfigName("config")
		viper.AutomaticEnv()
		viper.SetConfigType(".env")
		viper.AddConfigPath(".")

		config = &Config{
			Port: viper.GetInt("port"),
			Postgres: PostgresConfig{
				Host:          viper.GetString("db.postgres.host"),
				Port:          viper.GetInt("db.postgres.port"),
				Username:      viper.GetString("db.postgres.username"),
				Password:      viper.GetString("db.postgres.password"),
				Database:      viper.GetString("db.postgres.database"),
				MaxConnection: viper.GetInt("db.postgres.max_connection"),
			},
		}
	})

	return config
}
