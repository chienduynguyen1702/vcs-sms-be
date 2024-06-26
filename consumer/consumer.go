package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	uc_pb "vcs-sms-consumer/proto/uptime_calculate"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Consumer struct {
	DB          *gorm.DB
	KafkaReader *kafka.Reader
	ES          *ConsumerESClient
	Debug       bool
}

// Init is a method that initializes the Consumer
func InitConsumerInstance() *Consumer {
	// default ping interval is 300s
	h := &Consumer{}
	h.Debug = false
	return h
}

func (c *Consumer) SetDebugMode(debug bool) {
	c.Debug = debug
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
	// Buffer to store messages
	var messages []Server
	batchSize := 500
	for {

		// get 00h00m00s today and 23h59m59s today
		// start := time.Now().Truncate(24 * time.Hour)
		// end := start.Add(24*time.Hour - time.Second)
		// c.ES.AggregateUptimeServer(ES_INDEX_NAME, start, end)

		// get messages from kafka
		m, err := c.KafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Println("Failed to read message from kafka:", err)
			continue
		}

		server := Server{}
		err = server.UnmarshalJSON(m.Value)
		if err != nil {
			log.Println("Failed to unmarshal message from kafka:", err)
			continue
		}

		if c.Debug {
			server.PrintResult()
		}

		messages = append(messages, server)

		if len(messages) >= batchSize {
			err = c.bulkIndexToES(messages)
			if err != nil {
				log.Println("Failed to bulk index servers to elasticsearch:", err)
			}
			// Clear the messages slice after bulk indexing
			messages = messages[:0]
		}
	}
}

func (c *Consumer) bulkIndexToES(servers []Server) error {
	var buf bytes.Buffer

	for _, server := range servers {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "%s" } }%s`, ES_INDEX_NAME, "\n"))
		data, err := json.Marshal(server)
		if err != nil {
			return fmt.Errorf("failed to marshal server to JSON: %w", err)
		}
		buf.Write(meta)
		buf.Write(data)
		buf.WriteString("\n")
	}

	req := esapi.BulkRequest{
		Body: bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), c.ES.TypedClient)
	if err != nil {
		return fmt.Errorf("failed to execute bulk request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk request returned error: %s", res.Status())
	}

	return nil
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
	return nil
}

// set kafka reader
func (c *Consumer) SetKafkaReader(kafkaAddress *kafka.Reader) error {
	if kafkaAddress == nil {
		return fmt.Errorf("kafka reader is nil")
	}

	c.KafkaReader = kafkaAddress
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

// Query Elasticsearch for each server to aggregate online percentage of each server
func (c *Consumer) QueryES() {
	// query elasticsearch
	// get all servers
	// for each server, query elasticsearch
	// get the online percentage
	// update the server status
}

// RequestAggregation is a method that handles the RequestAggregation gRPC call
func (c *Consumer) RequestAggregation(ctx context.Context, req *uc_pb.AggregationRequest) (*uc_pb.AggregationResponse, error) {
	fd := req.GetFromDate()
	td := req.GetToDate()
	// parse fromDate and toDate to time.Time
	fromDate := fd.AsTime().Format(YYYYMMDD)
	toDate := td.AsTime().Format(YYYYMMDD)

	log.Printf("Client sent request to aggregation from %v to %v", fromDate, toDate)
	return &uc_pb.AggregationResponse{
		IsSuccess:                  true,
		AveragePercentUptimeServer: 100,
	}, nil
}
