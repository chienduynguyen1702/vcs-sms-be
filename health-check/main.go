package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

var dbCreds DBCredentials
var ctx = context.Background()
var initkafkaWriter *kafka.Writer
var err error
var interval int
var HealthCheckDebugMode bool

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
	initkafkaWriter, err = ConnectProducerToKafka(kafkaAddress, "ping_status")
	if err != nil {
		log.Println("Failed to connect to kafka:", err)
	}

	intervalStr := os.Getenv("PING_INTERVAL")
	//parse to int
	interval, err = strconv.Atoi(intervalStr)
	if err != nil {
		log.Println("Failed to parse interval:", err)
	}

	HealthCheckDebugMode, err = strconv.ParseBool(os.Getenv("HEALTH_CHECK_DEBUG_MODE"))
	if err != nil {
		log.Println("Failed to parse debug mode:", err, "setting debug mode to false")
		HealthCheckDebugMode = false
	}
}

// ################ main function ################
func main() {
	// create a new health check instance
	h := InitHealthCheckInstance()

	// connect to the database
	err := h.ConnectDB(dbCreds)
	if err != nil {
		panic("Failed to connect to database")
	}

	// connect to kafka
	err = h.SetKafkaWriter(initkafkaWriter)
	if err != nil {
		log.Println("Failed to connect to kafka:", err)

	}
	// set ping interval
	h.SetPingInterval(interval)

	// Set debug mode
	h.SetDebugMode(HealthCheckDebugMode)

	// validate the health check
	if !h.Validate() {
		panic("Failed to validate health check")
	}

	// start the health check
	h.StartHealthCheck()

}
