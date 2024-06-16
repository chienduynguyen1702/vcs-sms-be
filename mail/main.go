package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"vcs-sms-mail/proto/send_mail"
	"vcs-sms-mail/proto/uptime_calculate"

	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting mail service ...")
	if os.Getenv("RUN_ON_CONTAINER") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	mail := os.Getenv("MAIL_SERVICE_MAIL")
	pass := os.Getenv("MAIL_SERVICE_PASS")
	if mail == "" || pass == "" {
		fmt.Println("Please set MAIL_SERVICE_MAIL and MAIL_SERVICE_PASS")
		return
	}

	gmailService := InitGmailService(mail, pass)

	cron := cron.New()

	// @every day in 6am
	err := cron.AddFunc("0 0 0 * *", func() {
		fmt.Println("Starting send mail")

		err := gmailService.SendEmail("text.txt", "chiennd1702@gmail.com")
		if err != nil {
			fmt.Println("Error when send mail: ", err)
		}

		fmt.Println("Finish send mail")
	})

	if err != nil {
		fmt.Println("Error when add cron job: ", err)
		return
	}

	cron.Start()

	//################# GRPC SERVER #################
	grpcServer := grpc.NewServer()

	// Create UpTimeCalculateClient
	consumerAddress := os.Getenv("GRPC_SERVER_CONSUMER_ADDRESS")
	if consumerAddress == "" {
		fmt.Println("GRPC_SERVER_CONSUMER_ADDRESS is not set, using default address localhost:50051")
		consumerAddress = "localhost:50051"
	}
	// var opts []grpc.DialOption

	insecureCreds := insecure.NewCredentials()
	cc, err := grpc.NewClient(consumerAddress, grpc.WithTransportCredentials(insecureCreds))
	if err != nil {
		log.Println("Failed to create Client Con to Consumer server", err)
		panic(err)
	}
	uptimeClient := uptime_calculate.NewUptimeCalculateClient(cc)
	// log
	log.Printf("Connected to Consumer server at %s", consumerAddress)

	// Create SendMailServerImpl
	smServer := &SendMailServerImpl{gmail: gmailService, uptimeClient: uptimeClient}

	// start the gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		fmt.Println("GRPC_PORT is not set, using default port 50051")
		grpcPort = "50052"
	}
	grpcAddress := ":" + grpcPort

	// fmt.Println("grpcAddress:", grpcAddress)

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	send_mail.RegisterSendMailServer(grpcServer, smServer)
	log.Printf("Mail GRPC server listening at %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
