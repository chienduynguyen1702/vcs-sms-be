package main

import (
	"context"
	"log"
	"vcs-sms-mail/proto/send_mail"
)

type SendMailServerImpl struct {
	gmail *GmailService
	send_mail.UnimplementedSendMailServer
}

func (s *SendMailServerImpl) DoSendMail(ctx context.Context, in *send_mail.MailRequest) (*send_mail.MailResponse, error) {
	//debug
	// fmt.Println("MailReceiver: ", in.MailReceiver)
	// fmt.Println("FromDate: ", in.FromDate)
	// fmt.Println("ToDate: ", in.ToDate)
	// fmt.Println("TotalServer: ", in.TotalServer)
	// fmt.Println("NumberOfOnlineServer: ", in.NumberOfOnlineServer)
	// fmt.Println("NumberOfOfflineServer: ", in.NumberOfOfflineServer)
	// fmt.Println("AveragePercentUptimeServer: ", in.AveragePercentUptimeServer)

	// convert *timestamppb.Timestamp to time.Time
	fromDate := in.FromDate.AsTime()
	toDate := in.ToDate.AsTime()
	// Send mail
	log.Println("Sending mail to:", in.MailReceiver)

	err := s.gmail.SendEmailV2(fromDate, toDate, in.TotalServer, in.NumberOfOnlineServer, in.NumberOfOfflineServer, in.AveragePercentUptimeServer, in.MailReceiver)
	if err != nil {
		log.Println("Failed to send mail to:", in.MailReceiver)
		log.Println("Error:", err)
		return &send_mail.MailResponse{
			IsSuccess: false,
			Message:   err.Error(),
		}, err
	}

	log.Println("Sent mail to:", in.MailReceiver)
	log.Println("Reponse grpc")

	return &send_mail.MailResponse{
		IsSuccess: true,
		Message:   "Sent mail to: " + in.MailReceiver,
	}, nil
}
