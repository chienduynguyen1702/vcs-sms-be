package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"vcs-sms-consumer/proto/uptime_calculate"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
	YYYYMMDD       = "2006-01-02"
	ES_INDEX_NAME  = "server"
)

var dbCreds DBCredentials
var ctx = context.Background()
var initkafkaReader *kafka.Reader
var err error
var cloudID string
var apiKey string

// var esCloudAddress string

// define a gRPC server attibutes
var grpcServer *grpc.Server
var grpcPort string

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

	// create a new gRPC server
	grpcServer = grpc.NewServer()
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
	// ######## start the health check
	go c.StartConsumer()

	//################# GRPC SERVER #################
	// Create UpTimeCalculateServer
	ucServer := &UptimeCalculateServerImpl{consumer: c}

	// start the gRPC server
	grpcPort = os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		fmt.Println("GRPC_PORT is not set, using default port 50051")
		grpcPort = "50051"
	}
	grpcAddress := ":" + grpcPort

	// fmt.Println("grpcAddress:", grpcAddress)

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	uptime_calculate.RegisterUptimeCalculateServer(grpcServer, ucServer)
	log.Printf("Consumer GRPC Server listening at %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
