package main

import "gopkg.in/gomail.v2"

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

func (m *GmailService) SendEmail(filePath string, toMail string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "no-reply-vcs@mail.vn")
	msg.SetHeader("To", toMail)
	msg.SetHeader("Subject", "Statistical Servers Report")
	msg.SetBody("text/plain", "Here is the report you requested.")
	msg.Attach(filePath)

	if err := m.Dialer.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
