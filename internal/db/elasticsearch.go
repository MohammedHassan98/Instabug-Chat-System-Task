package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func setupElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
		Username:  os.Getenv("ELASTICSEARCH_USERNAME"),
		Password:  os.Getenv("ELASTICSEARCH_PASSWORD"),
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	ES = client

	// Create message index with mapping
	log.Println("Elasticsearch client created Successfully")
	createMessageIndex()
}

func createMessageIndex() error {
	mapping := `{
		"mappings": {
			"properties": {
				"chat_id": { "type": "keyword" },
				"message_number": { "type": "integer" },
				"body": {
					"type": "text",
					"analyzer": "standard"
				},
				"created_at": { "type": "date" }
			}
		}
	}`

	res, err := ES.Indices.Create(
		"messages",
		ES.Indices.Create.WithBody(strings.NewReader(mapping)),
	)

	if err != nil {
		return fmt.Errorf("cannot create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}
