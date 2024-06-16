package dtos

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
