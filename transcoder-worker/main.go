package main

import (
	"log"
	"net/http"

	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/communicator"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/consumer"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	communicator := communicator.NewCommunicator(config.VideoServiceAddress, config.Port)
	consumer := consumer.NewConsumer(config.KafkaHost, config.KafkaTopic, communicator)
	defer consumer.Close()

	router := gin.Default()
	router.StaticFS("/", http.Dir("uploads"))
	router.Run(config.Port)
}

type Config struct {
	KafkaHost           string `mapstructure:"KAFKA_HOST"`
	KafkaTopic          string `mapstructure:"KAFKA_TOPIC"`
	Port                string `mapstructure:"PORT"`
	VideoServiceAddress string `mapstructure:"VIDEO_SERVICE_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
