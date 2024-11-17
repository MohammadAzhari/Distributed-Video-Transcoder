package main

import (
	"context"
	"log"

	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/api"
	db "github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/db/sqlc"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/video-service/producer"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	kafkaHost = "localhost:9092"
	topic     = "video"
)

func main() {
	producer := producer.NewProducer(kafkaHost, topic)
	defer producer.Close()

	connPool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/dvt?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(producer, store)
	server.Start(":8080")
}
