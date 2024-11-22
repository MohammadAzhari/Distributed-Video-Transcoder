package main

import (
	"context"
	"log"

	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/api"
	db "github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/db/sqlc"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/producer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	producer := producer.NewProducer(config.KafkaHost, config.KafkaTopic)
	defer producer.Close()

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(producer, store)
	server.Start(config.Port)
}

type Config struct {
	DBSource   string `mapstructure:"DB_SOURCE"`
	KafkaHost  string `mapstructure:"KAFKA_HOST"`
	KafkaTopic string `mapstructure:"KAFKA_TOPIC"`
	Port       string `mapstructure:"PORT"`
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
