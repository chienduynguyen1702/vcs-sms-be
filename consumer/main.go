package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

var dbCreds DBCredentials
var ctx = context.Background()
var initkafkaReader *kafka.Reader
var err error

func init() {
	if os.Getenv("LOAD_ENV_FILE") != "true" {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}

	dbCreds = DBCredentials{
		Host: os.Getenv("DB_HOST"),
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Name: os.Getenv("DB_NAME"),
		Port: os.Getenv("DB_PORT"),
	}

	kafkaAddress := os.Getenv("KAFKA_BROKER")
	initkafkaReader = NewKafkaReader(pingStatusTopicName, kafkaAddress)
	if err != nil {
		log.Println("Failed to connect to kafka:", err)
	}
}

// ################ main function ################
func main() {
	// create a new health check instance
	c := InitConsumerInstance()

	// connect to the database
	err := c.ConnectDB(dbCreds)
	if err != nil {
		panic("Failed to connect to database")
	}

	// connect to kafka
	err = c.SetKafkaReader(initkafkaReader)
	if err != nil {
		log.Println("Failed to connect to kafka:", err)

	}

	// validate the health check
	if !c.Validate() {
		panic("Failed to validate health check")
	}

	// start the health check
	c.StartConsumer()

}
