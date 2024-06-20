package main

import (
	"fmt"
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
func (m *GmailService) SendEmail(toMail string) error {

	msg := gomail.NewMessage()
	msg.SetHeader("From", "chienduynguyen1702@mail.com")
	msg.SetHeader("To", toMail)
	msg.SetHeader("Subject", "Statistical Servers Report")
	msg.SetBody("text/plain", "Here is the report you requested.")

	if err := m.Dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func (m *GmailService) DefaultFilePath() string {
	todayString := time.Now().Format("2006-01-02")

	return DEFAULT_DIR + "/" + DEFAULT_PREFIX_FILE_NAME + todayString + ".xlsx"
}

// send file at filePath to toMail, default filePath is /reports/VCS-SMS-Report-<today>.xlsx
func (m *GmailService) SendEmailV2(FromDate, ToDate time.Time, TotalServer, NumberOfOnlineServer, NumberOfOfflineServer int64, AveragePercentUptimeServer float32, MailReceiver string) error {
	avgfloat64 := float64(AveragePercentUptimeServer)
	avgPercent := fmt.Sprintf("%.2f", avgfloat64*100) // convert to percent
	body := `
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 20px;
					color: #333;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					padding: 20px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
					border-radius: 8px;
				}
				h1 {
					color: #4CAF50;
				}
				p {
					font-size: 16px;
				}
				.data {
					background-color: #f9f9f9;
					padding: 10px;
					border-radius: 5px;
					margin-top: 10px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>[VCS-SMS] Statistical Servers Report</h1>
				<p>Total Server: <strong>` + fmt.Sprint(TotalServer) + `</strong></p>
				<p>Number of Online Server: <strong>` + fmt.Sprint(NumberOfOnlineServer) + `</strong></p>
				<p>Number of Offline Server: <strong>` + fmt.Sprint(NumberOfOfflineServer) + `</strong></p>
				<div class="data">
					<p>From: <strong>` + FromDate.Format("2006-01-02") + `</strong></p>
					<p>To: <strong>` + ToDate.Format("2006-01-02") + `</strong></p>
					<p>Average Percent Uptime Server: <strong>` + avgPercent + `%</strong></p>
				</div>
			</div>
		</body>
		</html>
	`
	msg := gomail.NewMessage()
	msg.SetHeader("From", "chienduynguyen1702@mail.com")
	msg.SetHeader("To", MailReceiver)
	msg.SetHeader("Subject", "Statistical Servers Report")
	msg.SetBody("text/html", body)

	if err := m.Dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
