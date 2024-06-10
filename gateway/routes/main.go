// routes/router.go

package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	docs "github.com/chienduynguyen1702/vcs-sms-be/docs"
	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestResponseLogger logs request and response details
func RequestResponseLogger() gin.HandlerFunc {
	// Cấu hình Lumberjack để quản lý log rotation
	logFile := &lumberjack.Logger{
		Filename:   "./logs/api.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	return func(c *gin.Context) {
		startTime := time.Now()

		// Đọc request body
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		c.Next()

		// Tính toán thời gian phản hồi
		latency := time.Since(startTime)

		// Trích xuất thông tin response
		responseBody := bodyLogWriter.body.String()
		// Bind Json responseBody to dtos.Response
		var response dtos.Response
		err := json.Unmarshal([]byte(responseBody), &response)
		if err != nil {
			response = dtos.Response{
				Success: false,
				Message: "Internal Server Error: Not found",
				Data:    nil,
			}
		}

		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		logEntry := fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s |%5d| %#v\n",
			startTime.Format("2006/01/02 - 15:04:05"),
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Writer.Size(),
			c.Request.URL.Path,
		)
		// Log response message if status code indicates an error
		if c.Writer.Status() >= 400 {
			logEntry = fmt.Sprintf("%s%s\n", logEntry, response.Message)
		}
		io.MultiWriter(logFile).Write([]byte(logEntry))
	}
}

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
	r.Use(gin.Recovery())
	// use logger
	r.Use(RequestResponseLogger())

	// Setup routes for the API version 1
	v1 := r.Group("/api/v1")
	{
		setupGroupUserRole(v1)
		setupGroupAuth(v1)
		setupGroupOrganization(v1)
		setupGroupServer(v1)
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
