// routes/router.go

package routes

import (
	"os"
	"time"

	docs "github.com/chienduynguyen1702/vcs-sms-be/docs"
	"github.com/chienduynguyen1702/vcs-sms-be/factory"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupV1Router() *gin.Engine {
	r := gin.Default()
	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("HOSTNAME_URL"), "http://localhost:" + os.Getenv("PORT")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 30 * time.Hour,
	}))

	// Setup routes for the API version 1
	v1 := r.Group("/api/v1")
	{
		setupGroupUser(v1)
		setupGroupAuth(v1)
		setupGroupOrganization(v1)
	}
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Swagger setup
	docs.SwaggerInfo.Title = "Parameter Store Backend API"
	docs.SwaggerInfo.Description = "This is a simple API for Parameter Store Backend."
	docs.SwaggerInfo.Version = "1.0"
	if os.Getenv("ENVIRONMENT") == "dev" {
		docs.SwaggerInfo.Host = "localhost:" + os.Getenv("PORT")
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	} else if os.Getenv("ENVIRONMENT") == "production" {
		docs.SwaggerInfo.Host = "vcs-sms-be-golang.up.railway.app"
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else if os.Getenv("ENVIRONMENT") == "datn-server" {
		docs.SwaggerInfo.Host = os.Getenv("HOSTNAME")
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	mainController := factory.AppFactoryInstance.CreateMainController()
	v1.GET("/ping", mainController.Ping)

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
