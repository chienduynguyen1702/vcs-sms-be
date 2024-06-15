package main

import (
	"context"
	"log"
	"time"
	uc_pb "vcs-sms-consumer/proto/uptime_calculate"
)

type UptimeCalculateServer interface {
	RequestAggregation(context.Context, *uc_pb.AggregationRequest) (*uc_pb.AggregationResponse, error)
}

// UptimeCalculateServerImpl implements UptimeCalculateServer
type UptimeCalculateServerImpl struct {
	consumer *Consumer
	uc_pb.UnimplementedUptimeCalculateServer
}

func (s *UptimeCalculateServerImpl) RequestAggregation(ctx context.Context, req *uc_pb.AggregationRequest) (*uc_pb.AggregationResponse, error) {
	fromDate := req.GetFromDate().AsTime()
	toDate := req.GetToDate().AsTime()

	// check if fromDate is after toDate
	if fromDate.After(toDate) {
		log.Println("fromDate is after toDate")
		return &uc_pb.AggregationResponse{
			IsSuccess: false,
			FilePath:  "",
		}, nil
	}

	// Get 00:00:00 of fromDate
	startDateOfFromDate := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
	// Get 23:59:59 of toDate
	endDateOfToDate := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 23, 59, 59, 999999999, toDate.Location())
	// Run Elasticsearch query
	filePath, err := s.consumer.ES.AggregateUptimeServer(ES_INDEX_NAME, startDateOfFromDate, endDateOfToDate)
	if err != nil {
		return &uc_pb.AggregationResponse{
			IsSuccess: false,
			FilePath:  "",
		}, err
	}

	// Process and return the response
	response := &uc_pb.AggregationResponse{
		// Populate the response with esResult data
		IsSuccess: true,
		FilePath:  filePath,
	}

	return response, nil
}
