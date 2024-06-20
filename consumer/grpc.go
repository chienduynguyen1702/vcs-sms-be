package main

import (
	"context"
	"log"
	"time"
	"vcs-sms-consumer/proto/uptime_calculate"
)

// UptimeCalculateServerImpl implements UptimeCalculateServer
type UptimeCalculateServerImpl struct {
	consumer *Consumer
	uptime_calculate.UnimplementedUptimeCalculateServer
}

func (s *UptimeCalculateServerImpl) RequestAggregation(ctx context.Context, req *uptime_calculate.AggregationRequest) (*uptime_calculate.AggregationResponse, error) {
	fromDate := req.GetFromDate().AsTime()
	toDate := req.GetToDate().AsTime()

	// check if fromDate is after toDate
	if fromDate.After(toDate) {
		log.Println("fromDate is after toDate")
		return &uptime_calculate.AggregationResponse{
			IsSuccess:                  false,
			AveragePercentUptimeServer: 0,
		}, nil
	}

	// Get 00:00:00 of fromDate
	startDateOfFromDate := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
	// Get 23:59:59 of toDate
	endDateOfToDate := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 23, 59, 59, 999999999, toDate.Location())
	// Run Elasticsearch query
	result, err := s.consumer.ES.AggregateUptimeServer(ES_INDEX_NAME, startDateOfFromDate, endDateOfToDate)
	// fmt.Println("result", result)
	if err != nil {
		return &uptime_calculate.AggregationResponse{
			IsSuccess:                  false,
			AveragePercentUptimeServer: 0,
		}, err
	}

	// Process and return the response
	response := &uptime_calculate.AggregationResponse{
		// Populate the response with esResult data
		IsSuccess:                  true,
		AveragePercentUptimeServer: result,
	}

	return response, nil
}
