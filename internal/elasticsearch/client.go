package elasticsearch

import (
	"context"
	"fmt"
	"log"
	"time"

	elastic "github.com/olivere/elastic/v7"
)

type Client struct {
	ESClient *elastic.Client
}

func NewElasticsearchClient(esURL string) (*Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(nil, "ELASTIC_ERROR: ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(nil, "ELASTIC_INFO: ", log.LstdFlags)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	info, code, err := client.Ping(esURL).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping Elasticsearch at %s: %w", esURL, err)
	}
	log.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	return &Client{ESClient: client}, nil
}
