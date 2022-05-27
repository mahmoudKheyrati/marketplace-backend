package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	Port int
	//RoundRobinTimeQuantum                    int
	//MultilevelFeedbackQueueLevelsTimeQuantum []int
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
		}
	})

	return config
}
