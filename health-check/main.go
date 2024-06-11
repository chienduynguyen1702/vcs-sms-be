package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var dbCreds DBCredentials
var ctx = context.Background()

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

}

// ################ main function ################
func main() {
	// create a new health check instance
	h := HealthCheck{}
	// connect to the database
	err := h.ConnectDB(dbCreds)
	if err != nil {
		log.Println("Failed to connect to database")
	}
	// get all the servers
	for {
		// ################################################
		servers, err := h.GetServers()
		if err != nil {
			log.Println("Failed to get servers")
		}
		// ################################################
		// Try to ping
		h.PingServers(servers)
		// ################################################
		// update data

		// sleep for 3 min
		time.Sleep(5 * time.Second)
	}

}
