package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	orderTopicName = "orders"
)

func NewConn(kafkaAddress string) *kafka.Conn {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	return conn
}

// Write message to Kafka topic
func NewKafkaWriter(topic string, kafkaAddress string) *kafka.Writer {
	writeTimeoutEnv := os.Getenv("KAFKA_BATCH_TIMEOUT")
	writeTimeout, err := strconv.ParseInt(writeTimeoutEnv, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	writeBatchTimeout := time.Duration(writeTimeout) * time.Nanosecond
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{kafkaAddress},
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: writeBatchTimeout,
	})
}

// var kafkaBroker = []string{os.Getenv("KAFKA_BROKER")}
var kafkaAddress = os.Getenv("KAFKA_BROKER")
var orderTopicProducer *kafka.Writer

// Connect to Kafka and check if working topic is exist
func ConnectProducerToKafka() *kafka.Writer {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, orderTopicName, 0)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	if err := CreateOrderTopicIfNotExists(conn); err != nil {
		panic(err.Error())
	}
	orderTopicProducer = NewKafkaWriter(orderTopicName, kafkaAddress)
	return orderTopicProducer
}

// CreateTopicIfNotExists create topic if not exists
func CreateOrderTopicIfNotExists(conn *kafka.Conn) error {
	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	for _, p := range partitions {
		if p.Topic == orderTopicName {
			log.Printf("Topic %s is already existed.", orderTopicName)
			return nil
		}
	}

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             orderTopicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("Topic %s created.", orderTopicName)
	return nil
}
