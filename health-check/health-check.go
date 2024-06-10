package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HealthCheck struct {
	DB *gorm.DB
}

type DBCredentials struct {
	Host string
	User string
	Pass string
	Name string
	Port string
}

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
		log.Println("Failed to close database")
	}
	db.Close()
	log.Println("Database connection closed")
}

// GetServers is a method that gets all the servers
func (h *HealthCheck) GetServers() ([]Server, error) {
	var servers []Server

	err := h.DB.Table("servers").Find(&servers)
	if err.Error != nil {
		log.Println("Failed to get servers")
		return nil, err.Error
	}

	return servers, nil
}

func (h *HealthCheck) GetServer(id int) (Server, error) {
	var server Server
	h.DB.Table("servers").Where("id = ?", id).Find(&server)
	return server, nil
}

func (h *HealthCheck) PingServers(ListServers []Server) {
	// create go routines to ping servers and wait group

	wg := sync.WaitGroup{}
	for _, server := range ListServers {
		// ping each server
		wg.Add(1)
		go h.Ping(&server, &wg)
		// if server is down, send an email
	}
	wg.Wait()
	fmt.Println("Finish Ping Servers")
}
func (h *HealthCheck) Ping(server *Server, wg *sync.WaitGroup) bool {
	attempt := 0
	// try 3 times attempt
	for attempt < 3 {
		// // real ping server
		// if h.ping(server) {
		// 	server.Status = "alive"
		// 	wg.Done()
		// 	return true
		// }

		// fake ping server
		if h.fakePing(server) {
			server.Status = "alive"
			server.PrintOne()
			wg.Done()
			return true
		}

		// if can't ping server, try again
		attempt++
	}
	server.Status = "dead"
	server.PrintOne()
	wg.Done()
	return false
}

// generate a fake ping 95% is alive
func (h *HealthCheck) fakePing(server *Server) bool {

	// 95% chance of being alive
	randomNumber := rand.Intn(100)
	time.Sleep(500 * time.Millisecond)
	fmt.Println(randomNumber)
	return randomNumber < 50
}

func (h *HealthCheck) ping(server *Server) bool {
	out, _ := exec.Command("ping", server.IP, "-w 500").Output()
	if strings.Contains(string(out), "bytes=") {
		return true
	} else {
		return false
	}
}
