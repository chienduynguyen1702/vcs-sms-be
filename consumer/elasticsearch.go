package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
)

type ConsumerESClient struct {
	TypedClient *elasticsearch.TypedClient
}

// create a new connection to elasticsearch
func NewTypedClientConnection(cloudID, apiKey string) (*elasticsearch.TypedClient, error) {

	// Es config
	cfg := elasticsearch.Config{
		// Addresses: []string{
		// 	address,
		// },
		// Username: os.Getenv("ELASTICSEARCH_USERNAME"),
		// Password: os.Getenv("ELASTICSEARCH_PASSWORD"),

		CloudID: cloudID,
		APIKey:  apiKey,
	}

	// Connect to Elasticsearch
	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch TypedClient: %s", err)
		return nil, err
	}
	fmt.Println("Created Elasticsearch TypedClient !")
	return es, nil
}

func InitConsumerESClient() *ConsumerESClient {
	h := &ConsumerESClient{}
	return h
}

// SetESClient sets the elasticsearch client and check if the index exists
func (c *ConsumerESClient) SetESClient(cloudID, apiKey string) error {
	c.TypedClient, err = NewTypedClientConnection(cloudID, apiKey)
	if err != nil {
		return err
	}
	if c.IndexExists("server") {
		log.Println("Index server exists")
	} else {
		log.Println("Index server does not exist, creating one")
		c.CreateIndexServer("server")
	}
	return nil
}

// CreateIndexServer creates an index named "server" in elasticsearch
func (c *ConsumerESClient) CreateIndexServer(indexName string) {
	// Index a document
	_, err := c.TypedClient.Indices.Create(indexName).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"ping_at":   types.NewDateNanosProperty(),
					"name":      types.NewTextProperty(),
					"status":    types.NewTextProperty(),
					"ip":        types.NewIpProperty(),
					"is_online": types.NewBooleanProperty(),
				},
			},
		}).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error creating the index: %s", err)
	}
}

// Check if the index exists
func (c *ConsumerESClient) IndexExists(indexName string) bool {
	// Check if the index exists
	exist, err := c.TypedClient.Indices.Exists(indexName).Do(context.Background())
	if err != nil {
		log.Fatalf("Error checking if the index exists: %s", err)
	}
	return exist
}

// IndexServer indexes a server in elasticsearch
func (c *ConsumerESClient) IndexServer(indexName string, server Server) error {
	document := map[string]interface{}{
		"ping_at":   server.PingAt,
		"ip":        server.IP,
		"status":    server.Status,
		"is_online": server.Status == "Online",
	}

	// Index a document
	_, err := c.TypedClient.Index(indexName).
		Request(document).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsumerESClient) AggregateUptimeServer(indexName string, startTime, toTime time.Time) {
	// build query
	agg := c.aggregationUptimeServerBuilder(startTime, toTime)

	// Perform the aggregation
	res, err := c.TypedClient.Search().
		Index(indexName).
		Request(agg).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error aggregating uptime: %s", err)
	}

	// Parse and extract the results
	// var result map[string]interface{}
	// if err := json.NewDecoder(res).Decode(&result); err != nil {
	// 	log.Fatalf("Error parsing the response body: %s", err)
	// }

	// unmarshall the response by below function
	// func (s *search.Response) UnmarshalJSON(data []byte) error
	// var resString string
	// res.UnmarshalJSON([]byte(res.Aggregations["date_filter"]))
	extractResults(res)
	// Print the results to check
}

const (
	DATE_FILTER_AGG = "date_filter"
	BY_DAY_AGG      = "by_day"
	BY_SERVER_AGG   = "by_server"
	ONLINE_PERCENT  = "online_percentage"
)

// aggregationUptimeServerBuilder builds the query for aggregating uptime of servers
func (c *ConsumerESClient) aggregationUptimeServerBuilder(startTime, toTime time.Time) *search.Request {
	ipField := "ip"
	pingAtField := "ping_at"
	isOnlineField := "is_online"
	startTimeStr := startTime.Format(time.RFC3339)
	toTimeStr := toTime.Format(time.RFC3339)
	searchRequest := &search.Request{
		Size: some.Int(0),
		Aggregations: map[string]types.Aggregations{ // Aggregate uptime of servers
			DATE_FILTER_AGG: {
				Filter: &types.Query{
					Range: map[string]types.RangeQuery{
						"ping_at": types.DateRangeQuery{
							Gte: &startTimeStr,
							Lte: &toTimeStr,
						},
					},
				},

				Aggregations: map[string]types.Aggregations{
					BY_DAY_AGG: {
						DateHistogram: &types.DateHistogramAggregation{
							Field:            &pingAtField,
							CalendarInterval: &calendarinterval.Day,
						},
						Aggregations: map[string]types.Aggregations{
							BY_SERVER_AGG: {
								Terms: &types.TermsAggregation{
									Field: &ipField,
									Size:  some.Int(10000),
								},
								Aggregations: map[string]types.Aggregations{
									"total_pings": {
										ValueCount: &types.ValueCountAggregation{
											Field: &isOnlineField,
										},
									},
									"online_pings": {
										Filter: &types.Query{
											Term: map[string]types.TermQuery{
												"is_online": {
													Value: true,
												},
											},
										},
										Aggregations: map[string]types.Aggregations{
											"online_count": {
												ValueCount: &types.ValueCountAggregation{
													Field: &isOnlineField,
												},
											},
										},
									},
									ONLINE_PERCENT: {
										BucketScript: &types.BucketScriptAggregation{
											BucketsPath: map[string]string{
												"total":  "total_pings",
												"online": "online_pings > online_count",
											},
											Script: &types.InlineScript{
												Source: "params.online / params.total * 100",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return searchRequest
}

type AggResult struct {
	Aggregations Aggregations `json:"aggregations"`
}
type Aggregations struct {
	DateFilter struct {
		DocCount int `json:"doc_count"`
		ByDay    struct {
			DayBuckets []struct {
				KeyAsString string `json:"key_as_string"`
				Key         int64  `json:"key"`
				DocCount    int    `json:"doc_count"`
				ByServer    struct {
					DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
					SumOtherDocCount        int `json:"sum_other_doc_count"`
					ServerBuckets           []struct {
						Key        string `json:"key"`
						DocCount   int    `json:"doc_count"`
						TotalPings struct {
							Value int `json:"value"`
						} `json:"total_pings"`
						OnlinePings struct {
							DocCount    int `json:"doc_count"`
							OnlineCount struct {
								Value int `json:"value"`
							} `json:"online_count"`
						} `json:"online_pings"`
						OnlinePercentage struct {
							Value float64 `json:"value"`
						} `json:"online_percentage"`
					} `json:"buckets"`
				} `json:"by_server"`
			} `json:"buckets"`
		} `json:"by_day"`
	} `json:"date_filter"`
}

func extractResults(res *search.Response) {
	resJson, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// fmt.Printf("resJson %s", string(resJson))
	fmt.Println()
	fmt.Println()
	fmt.Println()

	// fmt.Printf("resJson %s", string(resJson))
	var a AggResult

	err = json.Unmarshal(resJson, &a)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	fmt.Println("a.Aggregations.DateFilter.ByDay", a.Aggregations.DateFilter.ByDay)

	// dateFilterAgg, found := res.Aggregations[DATE_FILTER_AGG]
	// if !found {
	// 	log.Println(DATE_FILTER_AGG, " not found in response")
	// 	return
	// }
	// fmt.Println(DATE_FILTER_AGG, ":", dateFilterAgg)
	// fmt.Println("dateFilterAgg: %+v", dateFilterAgg)
	// fmt.Println("Found ", dateFilterAgg.[BY_DAY_AGG], " in response")
	// byDayAgg, found := dateFilterAgg.ChildrenAggregate
	// if !found {
	// 	log.Println("by_day not found in response")
	// 	return
	// }
	// byServerAgg, found := byDayAgg.Aggregations[BY_SERVER_AGG]

	// for _, dayBucket := range byDayAgg.Buckets {
	// 	date := dayBucket.KeyAsString

	// 	fmt.Printf("Date: %s\n", date)

	// 	byServerAgg, found := dayBucket.Aggregations["by_server"]
	// 	if !found {
	// 		log.Println("by_server not found in response")
	// 		continue
	// 	}

	// 	for _, serverBucket := range byServerAgg.Terms.Buckets {
	// 		ip := serverBucket.Key
	// 		onlinePercentage := serverBucket.Aggregations["online_percentage"].BucketScript.Value

	// 		fmt.Printf("Server IP: %s, Online Percentage: %.2f%%\n", *ip, *onlinePercentage)
	// 	}
	// }
}
