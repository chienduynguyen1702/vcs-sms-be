package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const MAXIMUM_CHANNELS = 100

type HealthCheck struct {
	DB          *gorm.DB
	KafkaWriter *kafka.Writer

	// PingInterval is the interval between pings, default is 300s
	PingInterval int
	Debug        bool

	GateWayAPIEndpoint string
}

// Init is a method that initializes the HealthCheck
func InitHealthCheckInstance() *HealthCheck {
	// default ping interval is 300s
	h := &HealthCheck{}
	h.PingInterval = 300
	h.Debug = false

	return h
}

// SetPingInterval is a method that sets the ping interval
func (h *HealthCheck) SetPingInterval(interval int) {
	h.PingInterval = interval
}

// SetDebug is a method that sets the debug mode
func (h *HealthCheck) SetDebugMode(debug bool) {
	h.Debug = debug
}

// SetGatewayAPIEndpoint is a method that sets the gateway API endpoint
func (h *HealthCheck) SetGatewayAPIEndpoint(endpoint string) {
	h.GateWayAPIEndpoint = endpoint
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
	if h.Debug {
		fmt.Println("| Debug        | Enabled   |")
	} else {
		fmt.Println("| Debug        | Disabled  |")
	}
	fmt.Printf("| PingInterval | %8ds |\n", h.PingInterval)
	if h.GateWayAPIEndpoint != "" {
		fmt.Printf("| GatewayAPI   | %s \n", h.GateWayAPIEndpoint)
	} else {
		fmt.Println("| GatewayAPI   | Not yet   ")
	}
	fmt.Println("----------------------------")
	fmt.Println("")
}

// StartHealthCheck is a method that starts the health check
func (h *HealthCheck) StartHealthCheck() {
	fmt.Println("")
	for {
		fmt.Println("")
		fmt.Println(" ========== Starting health check ")
		fmt.Println("")
		// get all servers
		servers, err := h.GetServers()
		if err != nil {
			log.Println("Failed to get servers")
		}
		// create go routines to ping servers and wait group
		wg := &sync.WaitGroup{}
		// ping servers
		h.PingServers(servers, wg)
		wg.Wait()

		fmt.Println("")
		h.SaveServers(servers)
		h.CallFlushCache(h.GateWayAPIEndpoint)

		fmt.Println(" ========== Finish health check ")
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
func (h *HealthCheck) GetServers() ([]*Server, error) {
	var servers []*Server

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
func (h *HealthCheck) PingServers(ListServers []*Server, wg *sync.WaitGroup) {

	// Create a buffered channel with capacity of MAXIMUM_CHANNELS
	channel := make(chan *Server, MAXIMUM_CHANNELS)

	// Function to process pinging each server
	pingServer := func(s *Server, wg *sync.WaitGroup) {
		err := h.Ping(s) // Assuming Ping method is defined on HealthCheck
		if err != nil {
			log.Println("Failed to ping server", s.IP, ":", err)
		}
		wg.Done()
	}

	// Start a goroutine to receive from the channel and ping servers
	go func() {
		for s := range channel {
			go pingServer(s, wg)
		}
	}()

	// Send servers to the channel
	for _, srv := range ListServers {
		wg.Add(1)
		channel <- srv
	}
	close(channel)

	// Wait for all pings to complete
	// wg.Wait()
}

func (h *HealthCheck) CallFlushCache(endpoint string) {
	// call api POST http://gateway:8080/api/v1/servers/flush-cache
	// using http client
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		log.Println("Failed to create request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to flush cache, status code:", resp.StatusCode)
	}
	defer resp.Body.Close()

	log.Println("Flush cache successfully")
}

// Save servers to database using transaction
func (h *HealthCheck) SaveServers(servers []*Server) {
	fmt.Println(" ==> Start save server to db")
	tx := h.DB.Begin()
	for index, server := range servers {
		// creata a save point
		savePointName := fmt.Sprintf("savepoint%d", index)
		tx.SavePoint(savePointName)
		err := tx.Save(&server).Error // save server to database
		if err != nil {
			tx.RollbackTo(savePointName)
			log.Println("Failed to save server", server.IP, ":", err)
			continue
		}
	}
	tx.Commit()
	fmt.Println(" ==> Done save server to db")
}

// This method pings all the servers then push results to kafka
func (h *HealthCheck) Ping(server *Server) error {
	// // ping server [REAL]
	// if h.ping(server) {
	// 	server.Status = "alive"
	// }

	// ping server [FAKE] : return 80% alive
	if h.fakePing() {
		server.Status = "Online"
		server.IsOnline = true
	} else {
		server.Status = "Offline"
		server.IsOnline = false
	}

	server.PingAt = time.Now()
	if h.Debug {
		server.PrintOne()
	}
	// push to kafka
	messageKafka := h.CreateMessage(*server)
	err := h.WriteMessageToKafka(messageKafka)
	if err != nil {
		log.Println("Failed to write message to kafka:", err)
		return err
	}
	return nil
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
