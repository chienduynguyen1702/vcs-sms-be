package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robfig/cron"

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
	err := cron.AddFunc("0 0 6 * *", func() {
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

	select {}
}
