package main

import (
	"fmt"
	"log"
	"os"

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
	// fmt.Println("Gmail service: ", gmailService.Mail)
	err := gmailService.SendEmail("text.txt", "chiennd1702@gmail.com")
	if err != nil {
		fmt.Println("Error when send mail: ", err)
	}
	fmt.Println("Finish send mail")
}
