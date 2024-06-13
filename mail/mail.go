package main

const (
	SUBJECT = "[VCS-SMS] Statistical Servers Report"
)

type Mail interface {
	SendEmail(filename string, toMail string) error
}

// func uploadHandler(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Save the file locally
// 	filename := "./" + file.Filename
// 	if err := c.SaveUploadedFile(file, filename); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Call the function to send email
// 	if err := sendEmail(filename); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and email sent successfully"})
// }

// func sendEmail(filename string) error {
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", "your-email@gmail.com")
// 	m.SetHeader("To", "recipient-email@gmail.com")
// 	m.SetHeader("Subject", "Excel File")
// 	m.SetBody("text/plain", "Here is the Excel file you requested.")
// 	m.Attach(filename)

// 	d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-email-password")

// 	if err := d.DialAndSend(m); err != nil {
// 		return err
// 	}
// 	return nil
// }
