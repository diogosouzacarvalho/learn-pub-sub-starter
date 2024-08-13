package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	fmt.Println()
	log.Println("Shutting down server...")
}
