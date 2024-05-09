package main

import (
	"log"
	"os"

	"github.com/chienduynguyen1702/vcs-sms-be/configs"
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/chienduynguyen1702/vcs-sms-be/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	if os.Getenv("SERVERLESS_DEPLOY") != "true" {
		configs.LoadEnvVariables()
	}
	db, err := configs.ConnectDatabase() // return *gorm.DB
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Migration db
	if os.Getenv("RUN_MIGRATION") == "true" {
		configs.Migration(db) // migration db
	}
	// Seed data

	// Set controller
	// controllers.SetDB(db) // set controller use that db *gorm.DB
	log.Println("Finished init.")
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	factory.AppFactoryInstance = factory.NewAppFactory(db)
}
func main() {

	port := os.Getenv("PORT")
	r := routes.SetupV1Router()
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	if os.Getenv("ENABLE_HTTPS_LOCAL") == "true" {
		log.Println("Server is running on port", port, "in", os.Getenv("GIN_MODE"), "gin mode with HTTPS")
		r.RunTLS(":"+port, os.Getenv("CERT_FILE_PATH"), os.Getenv("KEY_FILE_PATH"))
		return
	} else { // default
		log.Println("Server is running on port", port, "in", os.Getenv("GIN_MODE"), "gin mode")
		r.Run(":" + port)
	}
}
