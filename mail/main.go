package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"vcs-sms-mail/proto/send_mail"

	"github.com/robfig/cron"
	"google.golang.org/grpc"

	"github.com/joho/godotenv"
)

var gmailService *GmailService

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
	log.Printf("Mail: %s, Pass: %s", mail, pass)
	if mail == "" || pass == "" {
		fmt.Println("Please set MAIL_SERVICE_MAIL and MAIL_SERVICE_PASS")
		return
	}

	gmailService = InitGmailService(mail, pass)

	cron := cron.New()

	// @every day in 0am +7 timezone => 7am in Vietnam timezone
	//
	err := cron.AddFunc("0 0 0 * * *", func() {
		DailySendMail()
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

	// Create SendMailServerImpl
	smServer := &SendMailServerImpl{gmail: gmailService}

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

func DailySendMail() {
	var gatewayEndpoint string
	gatewayEndpoint = os.Getenv("GATEWAY_GET_MAIL_DATA_ENDPOINT")
	if gatewayEndpoint == "" {
		fmt.Println("GATEWAY_GET_MAIL_DATA_ENDPOINT is not set. Set as default http://gate-way:8080/api/v1/mail-infor")
		gatewayEndpoint = "http://gate-way:8080/api/v1/mail-infor"
		return
	}

	fmt.Println("Starting send mail")
	req, err := http.NewRequest("GET", gatewayEndpoint, nil)
	if err != nil {
		fmt.Println("Error when create request: ", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when send request: ", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error when send request, status code: ", resp.StatusCode)
		return
	}

	// fmt.Println("Response: ", resp)
	var mailReponse MailReponse
	err = json.NewDecoder(resp.Body).Decode(&mailReponse)
	if err != nil {
		fmt.Println("Error when decode response: ", err)
		return
	}
	mailBody := mailReponse.MailBody
	fmt.Println("MailReponse: ", mailBody)
	for _, mail := range mailBody.AdminMails {
		fmt.Println("Mail: ", mail)
		err := gmailService.SendEmailV2(mailBody.From, mailBody.To, mailBody.TotalServer, mailBody.TotalServerOnline, mailBody.TotalServerOffline, mailBody.AvgUptime, mail)
		if err != nil {
			fmt.Printf("Error when send mail to %s : %v", mail, err)
		}
	}

	fmt.Println("Finish send mail")
}
