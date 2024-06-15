package main

import (
	"context"
	"fmt"
	"log"
	"vcs-sms-mail/proto/send_mail"
	"vcs-sms-mail/proto/uptime_calculate"
)

type SendMailServerImpl struct {
	gmail *GmailService
	send_mail.UnimplementedSendMailServer
	uptimeClient uptime_calculate.UptimeCalculateClient
}

func (s *SendMailServerImpl) RequestAggregation(ctx context.Context, in *uptime_calculate.AggregationRequest) (*uptime_calculate.AggregationResponse, error) {
	// Call uptime calculate service
	uptimeResponse, err := s.uptimeClient.RequestAggregation(ctx, in)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Uptime response: ", uptimeResponse)

	return uptimeResponse, nil
}

func (s *SendMailServerImpl) DoSendMail(ctx context.Context, in *send_mail.MailRequest) (*send_mail.MailResponse, error) {

	//debug
	fmt.Println("MailReceiver: ", in.MailReceiver)
	fmt.Println("FromDate: ", in.FromDate)
	fmt.Println("ToDate: ", in.ToDate)

	// extract Date Range from request
	aggReq := &uptime_calculate.AggregationRequest{
		FromDate: in.FromDate,
		ToDate:   in.ToDate,
	}

	// Call uptime calculate service
	aggRes, err := s.RequestAggregation(ctx, aggReq)
	if err != nil {
		log.Println("Error when call uptime calculate service: ", err)
		return nil, err
	}

	// Bind MailResponse from AggregationResponse
	mailRes := &send_mail.MailResponse{
		IsSuccess: aggRes.IsSuccess,
		FilePath:  aggRes.FilePath,
	}

	// check if aggregation is not success
	if !aggRes.IsSuccess {
		log.Println("Aggregation is not success")
		return mailRes, nil
	}
	log.Println("Aggregation is success")
	log.Println("File path: ", aggRes.FilePath)

	log.Println("Sending mail to:", in.MailReceiver)
	// Send mail
	err = s.gmail.SendEmail(aggRes.FilePath, in.MailReceiver)
	if err != nil {
		log.Println("Error when send mail: ", err)
		return nil, err
	}
	log.Println("Sent mail to:", in.MailReceiver)
	return mailRes, nil
}
