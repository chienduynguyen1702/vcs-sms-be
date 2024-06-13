package main

import (
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Consumer struct {
	DB          *gorm.DB
	KafkaReader *kafka.Reader
	ES          *ConsumerESClient
}

// Init is a method that initializes the Consumer
func InitConsumerInstance() *Consumer {
	// default ping interval is 300s
	h := &Consumer{}

	return h
}

// Validate is a method that validates the Consumer struct
func (c *Consumer) Validate() bool {
	if c.DB == nil {
		log.Println("Database connection is nil, try ConnectDB() first")
		return false
	}
	if c.KafkaReader == nil {
		log.Println("Kafka reader is nil, try SetKafkaReader() first")
		return false
	}
	if c.ES == nil {
		log.Println("Elasticsearch client is nil, try SetESClient() first")
		return false
	}
	c.printValue()
	log.Println("Consumer is valid, ready to start")
	return true
}

// printValue is a method that prints the value of the Consumer struct
func (c *Consumer) printValue() {
	fmt.Println("----------------------------")
	fmt.Println("|      Consumer values     |")
	fmt.Println("----------------------------")

	if c.DB != nil {
		fmt.Println("| Database     | Connected |")
	} else {
		fmt.Println("| Database     | Not yet   |")
	}

	if c.KafkaReader != nil {
		fmt.Println("| KafkaReader  | Connected |")
	} else {
		fmt.Println("| KafkaReader  | Not yet   |")
	}

	if c.ES != nil {
		fmt.Println("| ES           | Connected |")
	} else {

		fmt.Println("| ES           | Not yet   |")
	}
	fmt.Println("----------------------------")
	fmt.Println("")
}

// StartConsumer is a method that starts the health check
func (c *Consumer) StartConsumer() {
	fmt.Println("")
	fmt.Println(" ========== Starting Consumer ==========")
	fmt.Println("")
	for {
		// get messages from kafka
		m, err := c.KafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Println("Failed to read message from kafka:", err)
		}
		// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// json unmarshall message to server struct
		server := Server{}
		err = server.UnmarshalJSON(m.Value)
		if err != nil {
			log.Println("Failed to unmarshal message from kafka:", err)
		}
		server.PrintResult()

	}
}

// use db to connect to database
func (c *Consumer) ConnectDB(dbCreds DBCredentials) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbCreds.Host,
		dbCreds.User,
		dbCreds.Pass,
		dbCreds.Name,
		dbCreds.Port,
	)
	// dsn : database source name
	// dsn := "host=localhost user=postgres password=postgres dbname=learning_gorm port=5432 sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database")
		return err
	}
	c.DB = DB
	// c.PingDB()
	log.Printf("Database connected\n")
	return nil
}

// set kafka reader
func (c *Consumer) SetKafkaReader(kafkaAddress *kafka.Reader) error {
	if kafkaAddress == nil {
		return fmt.Errorf("kafka reader is nil")
	}

	c.KafkaReader = kafkaAddress
	log.Println("Kafka reader connected")
	return nil
}

// pingDB is a method that pings the database
func (c *Consumer) PingDB() {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			log.Println("Failed to ping database")
		}
		err = sqlDB.Ping()
		if err != nil {
			log.Println("Failed to ping database")
		}
		log.Println("Database pinged")
	}
}

// CloseDB is a method that closes the database connection
func (c *Consumer) CloseDB() {
	db, err := c.DB.DB()
	if err != nil {
		log.Println("Failed to close database:", err)
	}
	db.Close()
	log.Println("Database connection closed")
}
