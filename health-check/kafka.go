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
	pingStatusTopicName = "ping_status"
)

func NewConn(kafkaAddress string) *kafka.Conn {
	// fmt.Println("Kafka address: ", kafkaAddress)
	conn, err := kafka.Dial("tcp", kafkaAddress)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	return conn
}

// Write message to Kafka topic
func newKafkaWriter(topic string, kafkaAddress string) *kafka.Writer {
	writeTimeoutEnv := os.Getenv("KAFKA_BATCH_TIMEOUT")
	writeTimeout, err := strconv.ParseInt(writeTimeoutEnv, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	writeBatchTimeout := time.Duration(writeTimeout) * time.Second
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{kafkaAddress},
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: writeBatchTimeout,
	})
}

// var kafkaBroker = []string{os.Getenv("KAFKA_BROKER")}
// var kafkaAddress = os.Getenv("KAFKA_BROKER")
// var pingStatusTopicProducer *kafka.Writer

// Connect to Kafka and check if working topic is exist
func ConnectProducerToKafka(kafkaAddress, pingStatusTopicName string) (*kafka.Writer, error) {
	// log.Printf("kafkaAddress: %s", kafkaAddress)
	// log.Printf("pingStatusTopicName: %s", pingStatusTopicName)

	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, pingStatusTopicName, 0)
	if err != nil {
		log.Printf("Failed to DialLeader to Kafka: %s", err)
		return nil, err
	}
	defer conn.Close()

	if err := createPingStatusTopicIfNotExists(conn); err != nil {
		return nil, err
	}
	pingStatusTopicProducer := newKafkaWriter(pingStatusTopicName, kafkaAddress)
	// log.Printf("pingStatusTopicProducer: %v", pingStatusTopicProducer)
	return pingStatusTopicProducer, nil
}

// CreateTopicIfNotExists create topic if not exists
func createPingStatusTopicIfNotExists(conn *kafka.Conn) error {
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}

	for _, p := range partitions {
		if p.Topic == pingStatusTopicName {
			log.Printf("Topic %s is already existed.", pingStatusTopicName)
			return nil
		}
	}

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             pingStatusTopicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		return err
	}
	log.Printf("Topic %s created.", pingStatusTopicName)
	return nil
}
