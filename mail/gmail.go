package main

import (
	"time"

	"gopkg.in/gomail.v2"
)

const (
	DEFAULT_DIR              = "/reports"
	DEFAULT_PREFIX_FILE_NAME = "VCS-SMS-Report-"
)

type GmailService struct {
	Dialer *gomail.Dialer
	Mail   string
	Pass   string
}

func InitGmailService(mail, pass string) *GmailService {
	d := gomail.NewDialer("smtp.gmail.com", 587, mail, pass)
	return &GmailService{
		Dialer: d,
		Mail:   mail,
		Pass:   pass,
	}
}

// send file at filePath to toMail, default filePath is /reports/VCS-SMS-Report-<today>.xlsx
func (m *GmailService) SendEmail(filePath string, toMail string) error {

	if filePath == "" {
		filePath = m.DefaultFilePath()
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", "chienduynguyen1702@mail.com")
	msg.SetHeader("To", toMail)
	msg.SetHeader("Subject", "Statistical Servers Report")
	msg.SetBody("text/plain", "Here is the report you requested.")
	msg.Attach(filePath)

	if err := m.Dialer.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

func (m *GmailService) DefaultFilePath() string {
	todayString := time.Now().Format("2006-01-02")

	return DEFAULT_DIR + "/" + DEFAULT_PREFIX_FILE_NAME + todayString + ".xlsx"
}
