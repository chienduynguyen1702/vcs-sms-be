package main

import (
	"context"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
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
func (c *ConsumerESClient) IndexServer(indexName string, server Server) {
	document := map[string]interface{}{
		"ping_at": server.PingAt,
		"name":    server.Name,
		"ip":      server.IP,
		"status":  server.Status,
	}

	// Index a document
	_, err := c.TypedClient.Index(indexName).
		Request(document).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error indexing the document: %s", err)
	}
}
