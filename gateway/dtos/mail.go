package dtos

import "fmt"

type SendMailRequest struct {
	Mail string `json:"mail"`
	From string `json:"from"` // from date
	To   string `json:"to"`   // to date
}

/*
{
  "from": "2024-06-13",
  "mail": "chiennd1702@gmail.com",
  "to": "2024-06-15"
}
*/

type MailBody struct {
	AdminMails         []string `json:"admin_mails"`
	From               string   `json:"from"`
	To                 string   `json:"to"`
	TotalServer        int64    `json:"total_server"`
	TotalServerOnline  int64    `json:"total_server_online"`
	TotalServerOffline int64    `json:"total_server_offline"`
	AvgUptime          float64  `json:"avg_uptime"`
}

func (m *MailBody) PrintMailBody() {
	fmt.Println("AdminMails: ", m.AdminMails)
	fmt.Println("From: ", m.From)
	fmt.Println("To: ", m.To)
	fmt.Println("TotalServer: ", m.TotalServer)
	fmt.Println("TotalServerOnline: ", m.TotalServerOnline)
	fmt.Println("TotalServerOffline: ", m.TotalServerOffline)
	fmt.Println("AvgUptime: ", m.AvgUptime)

}
