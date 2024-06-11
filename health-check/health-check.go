package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HealthCheck struct {
	DB          *gorm.DB
	KafkaWriter *kafka.Writer

	// PingInterval is the interval between pings, default is 300s
	PingInterval int
}

// Init is a method that initializes the HealthCheck
func InitHealthCheckInstance() *HealthCheck {
	// default ping interval is 300s
	h := &HealthCheck{}
	h.PingInterval = 300

	return h
}

// SetPingInterval is a method that sets the ping interval
func (h *HealthCheck) SetPingInterval(interval int) {
	h.PingInterval = interval
}

// Validate is a method that validates the HealthCheck struct
func (h *HealthCheck) Validate() bool {
	if h.DB == nil {
		log.Println("Database connection is nil, try ConnectDB() first")
		return false
	}
	if h.KafkaWriter == nil {
		log.Println("Kafka writer is nil, try SetKafkaWriter() first")
		return false
	}
	h.printValue()
	log.Println("HealthCheck is valid, ready to start")
	return true
}

// printValue is a method that prints the value of the HealthCheck struct
func (h *HealthCheck) printValue() {
	fmt.Println("----------------------------")
	fmt.Println("|    HealthCheck values    |")
	fmt.Println("----------------------------")
	if h.DB != nil {
		fmt.Println("| Database     | Connected |")
	} else {
		fmt.Println("| Database     | Not yet   |")
	}
	if h.KafkaWriter != nil {
		fmt.Println("| KafkaWriter  | Connected |")
	} else {
		fmt.Println("| KafkaWriter  | Not yet   |")
	}
	fmt.Printf("| PingInterval | %8ds |\n", h.PingInterval)
	fmt.Println("----------------------------")
	fmt.Println("")
}

// StartHealthCheck is a method that starts the health check
func (h *HealthCheck) StartHealthCheck() {
	fmt.Println("")
	for {
		fmt.Println(" ========== Starting health check ==========")
		fmt.Println("")
		// get all servers
		servers, err := h.GetServers()
		if err != nil {
			log.Println("Failed to get servers")
		}

		// ping servers
		h.PingServers(servers)

		fmt.Println("")
		fmt.Println(" ========== Finish health check ==========")

		// sleep interval time before pinging again
		time.Sleep(time.Duration(h.PingInterval) * time.Second)
	}
}

func (h *HealthCheck) createMessageString(serverResult Server) string {
	messageStr, err := json.Marshal(serverResult)
	if err != nil {
		log.Println("Failed to create message:", err)
		return ""
	}
	return string(messageStr)
}

// use Server status to create Kafka message
func (h *HealthCheck) CreateMessage(server Server) kafka.Message {
	messageStr := h.createMessageString(server)
	return kafka.Message{
		Key:   []byte("key"),
		Value: []byte(messageStr),
	}
}

// use db to connect to database
func (h *HealthCheck) ConnectDB(dbCreds DBCredentials) error {
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
	h.DB = DB
	// h.PingDB()
	log.Printf("Database connected\n")
	return nil
}

// set kafka writer
func (h *HealthCheck) SetKafkaWriter(kafkaAddress *kafka.Writer) error {
	if kafkaAddress == nil {
		return fmt.Errorf("kafka writer is nil")
	}

	h.KafkaWriter = kafkaAddress
	log.Println("Kafka writer connected")
	return nil
}

// pingDB is a method that pings the database
func (h *HealthCheck) PingDB() {
	if h.DB != nil {
		sqlDB, err := h.DB.DB()
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
func (h *HealthCheck) CloseDB() {
	db, err := h.DB.DB()
	if err != nil {
		log.Println("Failed to close database:", err)
	}
	db.Close()
	log.Println("Database connection closed")
}

// GetServers is a method that gets all the servers
func (h *HealthCheck) GetServers() ([]Server, error) {
	var servers []Server

	err := h.DB.Table("servers").Find(&servers)
	if err.Error != nil {
		log.Println("Failed to get servers:", err.Error)
		return nil, err.Error
	}

	return servers, nil
}

func (h *HealthCheck) GetServer(id int) (Server, error) {
	var server Server
	h.DB.Table("servers").Where("id = ?", id).Find(&server)
	return server, nil
}

// This method pings all the servers then push results to kafka
func (h *HealthCheck) PingServers(ListServers []Server) {
	// create go routines to ping servers and wait group

	wg := sync.WaitGroup{}
	for _, server := range ListServers {
		// ping each server
		wg.Add(1)
		go h.Ping(&server, &wg)

		time.Sleep(time.Duration(h.PingInterval) * time.Millisecond)
	}
	wg.Wait()
}

// This method pings all the servers then push results to kafka
func (h *HealthCheck) Ping(server *Server, wg *sync.WaitGroup) {
	// // ping server [REAL]
	// if h.ping(server) {
	// 	server.Status = "alive"
	// }

	// ping server [FAKE] : return 80% alive
	if h.fakePing() {
		server.Status = "Online"
	} else {
		server.Status = "Offline"
	}

	// if can't ping server, try again
	server.PrintOne()
	server.PingAt = time.Now()

	// push to kafka
	messageKafka := h.CreateMessage(*server)
	err := h.WriteMessageToKafka(messageKafka)
	if err != nil {
		log.Println("Failed to write message to kafka:", err)
	}

	// decrease wait group
	wg.Done()
}

// generate a fake ping 95% is alive
func (h *HealthCheck) fakePing() bool {

	// xx% chance of being alive
	randomNumber := rand.Intn(100)
	time.Sleep(500 * time.Millisecond)
	// fmt.Println(randomNumber)
	return randomNumber < 90
}

func (h *HealthCheck) ping(server *Server) bool {
	attempt := 0
	// try 3 times attempt
	for attempt < 3 {
		out, _ := exec.Command("ping", server.IP, "-w 500").Output()
		if strings.Contains(string(out), "bytes=") {
			return true
		}
		attempt++
	}

	return false
}

// Write message to kafka
func (h *HealthCheck) WriteMessageToKafka(message kafka.Message) error {
	// Write message to kafka
	if err := h.KafkaWriter.WriteMessages(ctx, message); err != nil {
		log.Panic("Failed to write message to kafka:", err)
		return err
	}
	return nil
}
