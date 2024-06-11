package main

import (
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	pingStatusTopicName = "ping_status"
)

// Reader of Kafka topic
func NewKafkaReader(topic, kafkaAddress string) *kafka.Reader {
	readTimeoutEnv := os.Getenv("KAFKA_READ_TIMEOUT")
	readTimeout, err := strconv.ParseInt(readTimeoutEnv, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	readBatchTimeout := time.Duration(readTimeout) * time.Nanosecond
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaAddress},
		GroupID:   "consumer-group-id",
		Topic:     pingStatusTopicName,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   readBatchTimeout,
	})
}
