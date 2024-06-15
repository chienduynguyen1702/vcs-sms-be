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
func NewKafkaReader(topic, kafkaAddress string) (*kafka.Reader, error) {
	readTimeoutEnv := os.Getenv("KAFKA_BATCH_TIMEOUT")
	readTimeout, err := strconv.ParseInt(readTimeoutEnv, 10, 64)
	if err != nil {
		return nil, err
	}
	readBatchTimeout := time.Duration(readTimeout) * time.Second
	// fmt.Println("Read batch timeout:", readBatchTimeout)
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaAddress},
		GroupID:   "consumer-group-id",
		Topic:     pingStatusTopicName,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   readBatchTimeout,
	}), nil
}
