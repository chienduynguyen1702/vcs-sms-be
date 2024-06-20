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

// BulkIndexServers indexes a batch of servers in elasticsearch
// func (c *ConsumerESClient) BulkIndexServers(indexName string, servers []Server) error {
// 	var bulkRequest []byte
// 	for _, server := range servers {
// 		// Create bulk request body
// 		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "%s" } }%s`, ES_INDEX_NAME, "\n"))
// 		data, err := json.Marshal(server)
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal server to JSON: %w", err)
// 		}
// 		bulkRequest = append(bulkRequest, meta...)
// 		bulkRequest = append(bulkRequest, data...)
// 		bulkRequest = append(bulkRequest, "\n"...)
// 	}

// 	// Perform bulk indexing
// 	err := c.TypedClient.API.Bulk(bytes.NewReader(bulkRequest))
// 	if err != nil {
// 		return fmt.Errorf("failed to execute bulk request: %w", err)
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		return fmt.Errorf("bulk request returned error: %s", res.Status())
// 	}

// 	return nil
// }

func (c *ConsumerESClient) AggregateUptimeServer(indexName string, startTime, toTime time.Time) (float32, error) {
	// build query
	agg := c.aggregationAvgUptimeServersBuilder(startTime, toTime)

	// Perform the aggregation
	res, err := c.TypedClient.Search().
		Index(indexName).
		Request(agg).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error aggregating uptime: %s", err)
		return -1, err
	}

	// Parse and extract the results
	avgResult := c.extractAvgResults(res)

	// debug
	// fmt.Println("Average uptime of servers from", startTime, "to", toTime, "is", avgResult)
	// parse float64 to float32
	avgResult32 := float32(avgResult)

	// Write to excel
	return avgResult32, nil
}

const (
	DATE_FILTER_AGG   = "date_filter"
	BY_DAY_AGG        = "by_day"
	BY_SERVER_AGG     = "by_server"
	ONLINE_PERCENT    = "online_percentage"
	OUTPUT_EXCEL_PATH = "/reports"
)

// aggregationUptimeServerBuilder builds the query for aggregating uptime of servers
func (c *ConsumerESClient) aggregationAvgUptimeServersBuilder(startTime, toTime time.Time) *search.Request {
	// ipField := "ip"
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
						pingAtField: types.DateRangeQuery{
							Gte: &startTimeStr,
							Lte: &toTimeStr,
						},
					},
				},

				Aggregations: map[string]types.Aggregations{
					"uptime_percent": {
						Avg: &types.AverageAggregation{
							Field: &isOnlineField,
						},
					},
				},
			},
		},
	}
	return searchRequest
}

// ########################################## Get average uptime of all servers fromdate-todate
// ElasticsearchResponse represents the full response from Elasticsearch
type AvgResponse struct {
	Took     int          `json:"took"`
	TimedOut bool         `json:"timed_out"`
	Shards   Shards       `json:"_shards"`
	Hits     Hits         `json:"hits"`
	Aggs     Aggregations `json:"aggregations"`
}

// Shards represents the shard information in the response
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// Hits represents the hits information in the response
type Hits struct {
	Total    TotalHits     `json:"total"`
	MaxScore interface{}   `json:"max_score"`
	Hits     []interface{} `json:"hits"`
}

// TotalHits represents the total hits information in the response
type TotalHits struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// Aggregations đại diện cho các kết quả tổng hợp trong phản hồi
type Aggregations struct {
	DateFilter DateFilter `json:"date_filter"`
}

// DateFilter đại diện cho kết quả tổng hợp date_filter
type DateFilter struct {
	DocCount      int           `json:"doc_count"`
	UptimePercent UptimePercent `json:"uptime_percent"`
}

// UptimePercent đại diện cho kết quả tổng hợp uptime_percent
type UptimePercent struct {
	Value         float64 `json:"value"`
	ValueAsString string  `json:"value_as_string"`
}

// get avg result
func (c *ConsumerESClient) extractAvgResults(res *search.Response) float64 {
	resJson, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// fmt.Printf("resJson %s", string(resJson))
	fmt.Println()
	fmt.Println()
	fmt.Println()

	// fmt.Printf("resJson %s", string(resJson))
	var a AvgResponse

	err = json.Unmarshal(resJson, &a)
	// fmt.Printf("AvgResponse %v", a)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	result := a.Aggs.DateFilter.UptimePercent.Value
	// fmt.Println("Result : ", result)

	return result
}

// ########################################## Get the UpTime percent of each server in each day fromdate-todate
// type AggResult struct {
// 	Aggregations Aggregations `json:"aggregations"`
// }

// type Aggregations struct {
// 	DateFilter struct {
// 		DocCount int `json:"doc_count"`
// 		ByDay    struct {
// 			DayBuckets []DayBuckets `json:"buckets"`
// 		} `json:"by_day"` // const BY_DAY_AGG
// 	} `json:"date_filter"` // const DATE_FILTER_AGG
// }

// type DayBuckets struct {
// 	KeyAsString string `json:"key_as_string"`
// 	Key         int64  `json:"key"`
// 	DocCount    int    `json:"doc_count"`
// 	ByServer    struct {
// 		DocCountErrorUpperBound int             `json:"doc_count_error_upper_bound"`
// 		SumOtherDocCount        int             `json:"sum_other_doc_count"`
// 		ServerBuckets           []ServerBuckets `json:"buckets"`
// 	} `json:"by_server"` // const BY_SERVER_AGG
// }

// func (d *DayBuckets) PrintResult() {
// 	fmt.Println("Day:", d.KeyAsString)
// 	for _, serverBucket := range d.ByServer.ServerBuckets {
// 		fmt.Println("Server:", serverBucket.Key)
// 		fmt.Println("Total pings:", serverBucket.TotalPings.Value)
// 		fmt.Println("Online pings:", serverBucket.OnlinePings.OnlineCount.Value)
// 		fmt.Println("Online percentage:", serverBucket.OnlinePercentage.Value)
// 	}
// 	fmt.Println("")
// }
// func PrintSliceOfDayBuckets(dayBuckets []DayBuckets) {
// 	for _, dayBucket := range dayBuckets {
// 		fmt.Println("Day:", dayBucket.KeyAsString)
// 		for _, serverBucket := range dayBucket.ByServer.ServerBuckets {
// 			fmt.Println("Server:", serverBucket.Key)
// 			fmt.Println("Total pings:", serverBucket.TotalPings.Value)
// 			fmt.Println("Online pings:", serverBucket.OnlinePings.OnlineCount.Value)
// 			fmt.Println("Online percentage:", serverBucket.OnlinePercentage.Value)
// 		}
// 		fmt.Println("")
// 	}
// }

// type ServerBuckets struct {
// 	Key        string `json:"key"`
// 	DocCount   int    `json:"doc_count"`
// 	TotalPings struct {
// 		Value int `json:"value"`
// 	} `json:"total_pings"`
// 	OnlinePings struct {
// 		DocCount    int `json:"doc_count"`
// 		OnlineCount struct {
// 			Value int `json:"value"`
// 		} `json:"online_count"`
// 	} `json:"online_pings"`
// 	OnlinePercentage struct {
// 		Value float64 `json:"value"`
// 	} `json:"online_percentage"` // const ONLINE_PERCENT
// }

// func (c *ConsumerESClient) extractResultsToDaily(res *search.Response) []DayBuckets {
// 	resJson, err := json.Marshal(res)
// 	if err != nil {
// 		log.Fatalf("Error parsing the response body: %s", err)
// 	}
// 	// fmt.Printf("resJson %s", string(resJson))
// 	fmt.Println()
// 	fmt.Println()
// 	fmt.Println()

// 	// fmt.Printf("resJson %s", string(resJson))
// 	var a AggResult

// 	err = json.Unmarshal(resJson, &a)
// 	if err != nil {
// 		log.Fatalf("Error parsing the response body: %s", err)
// 	}

// 	fmt.Println("Found", len(a.Aggregations.DateFilter.ByDay.DayBuckets), "daily aggregation buckets")

// 	return a.Aggregations.DateFilter.ByDay.DayBuckets
// }

// func (c *ConsumerESClient) WriteToExcel(dayBuckets []DayBuckets, startTime, toTime time.Time) string {
// 	fmt.Println("Writing to excel")
// 	// write to excel
// 	filePath := fmt.Sprintf("%s/VCS-SMS-Report-%s-%s.xlsx", OUTPUT_EXCEL_PATH, startTime.Format("2006-01-02"), toTime.Format("2006-01-02"))

// 	newFile := excelize.NewFile()
// 	// Create a new sheet by each DayBucket.KeyAsString
// 	for _, dayBucket := range dayBuckets {
// 		// Create a new sheet.
// 		sheetName := dayBucket.KeyAsString            // 2024-06-15T00:00:00.000Z
// 		t, err := time.Parse(time.RFC3339, sheetName) // convert to date format
// 		if err != nil {
// 			log.Println("Failed to parse time:", err)
// 			break
// 		}
// 		sheetName = t.Format("2006-01-02")

// 		index, err := newFile.NewSheet(sheetName)
// 		if err != nil {
// 			log.Println("Failed to create a new sheet: ", err)
// 			break
// 		}
// 		// Set value of a cell.
// 		newFile.SetCellValue(sheetName, "A1", "Server")
// 		newFile.SetCellValue(sheetName, "B1", "Online Percentage")

// 		for i, serverBucket := range dayBucket.ByServer.ServerBuckets {
// 			// Set value of Server column.
// 			err = newFile.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), serverBucket.Key)
// 			if err != nil {
// 				log.Println("Failed to et value of Server column:", err)
// 				break
// 			}
// 			// Set value of Online Percentage column.
// 			err = newFile.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), serverBucket.OnlinePercentage.Value)
// 			if err != nil {
// 				log.Println("Failed to Set value of Online Percentage column:", err)
// 				break
// 			}

// 		}
// 		// Set active sheet of the workbook.
// 		newFile.SetActiveSheet(index)
// 	}

// 	// Save xlsx file by the given path.
// 	if err := newFile.SaveAs(filePath); err != nil {
// 		log.Println("Failed to save xlsx:", err)
// 	}
// 	fmt.Println("Saved to excel file:", filePath)
// 	return filePath
// }
