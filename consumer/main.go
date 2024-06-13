package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
	ES_INDEX_NAME  = "server"
)

var dbCreds DBCredentials
var ctx = context.Background()
var initkafkaReader *kafka.Reader
var err error
var cloudID string
var apiKey string
var esCloudAddress string

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
	initkafkaReader, err = NewKafkaReader(pingStatusTopicName, kafkaAddress)
	if err != nil {
		log.Fatal("Failed to connect to kafka:", err)
	}
	cloudID = os.Getenv("ELASTICSEARCH_CLOUD_ID")
	apiKey = os.Getenv("ELASTICSEARCH_API_KEY")
	// esCloudAddress = os.Getenv("ELASTICSEARCH_CLOUD_ADDRESS")
	if cloudID == "" || apiKey == "" {
		log.Fatal("Failed to get elasticsearch credentials")
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
	log.Printf("Database connected\n")

	// connect to kafka
	err = c.SetKafkaReader(initkafkaReader)
	if err != nil {
		log.Println("Failed to connect to kafka:", err)
	}
	log.Println("Set KafkaReader !")
	c.ES = InitConsumerESClient()
	// connect to elasticsearch
	err = c.ES.SetESClient(cloudID, apiKey)
	if err != nil {
		log.Println("Failed to connect to elasticsearch:", err)
	}
	log.Println("Set ESClient !")

	// validate the health check
	if !c.Validate() {
		panic("Failed to validate health check")
	}
	c.SetDebugMode(true)
	// start the health check
	c.StartConsumer()

}
