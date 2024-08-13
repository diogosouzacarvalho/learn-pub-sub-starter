package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Starting Peril server...")

	connStr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBIT_MQ_SERVER_USER"),
		os.Getenv("RABBIT_MQ_SERVER_PASSWORD"),
		os.Getenv("RABBIT_MQ_SERVER_URL"),
		os.Getenv("RABBIT_MQ_SERVER_PORT"))

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Println(connStr)
		log.Fatalf("Error dialing rabbitmq server: %s", err)
	}
	defer conn.Close()
	log.Println("Connected to server!")

	if err := startGame(conn); err != nil {
		conn.Close()
		log.Fatalf("Error starting game: %s", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	fmt.Println()
	log.Println("Shutting down server...")
}

func startGame(c *amqp.Connection) (err error) {
	var channel *amqp.Channel

	if channel, err = c.Channel(); err != nil {
		return err
	}

	return pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
}
