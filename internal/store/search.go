package store

import (
	"context"
	"encoding/json"
	"fmt"

	"external-backend-go/internal/elasticsearch"
	"external-backend-go/internal/logger"

	elastic "github.com/olivere/elastic/v7"
)

type SearchStore interface {
	IndexDocument(ctx context.Context, indexName string, docID string, document interface{}) error
	Search(ctx context.Context, indexName string, query string, fields []string, page, pageSize int) ([]json.RawMessage, int64, error)
	DeleteDocument(ctx context.Context, indexName string, docID string) error
}

type genericSearchStore struct {
	esClient *elasticsearch.Client
	logger   *logger.Logger
}

func NewSearchStore(esClient *elasticsearch.Client, appLogger *logger.Logger) SearchStore {
	return &genericSearchStore{
		esClient: esClient,
		logger:   appLogger,
	}
}

func (s *genericSearchStore) IndexDocument(ctx context.Context, indexName string, docID string, document interface{}) error {
	_, err := s.esClient.ESClient.Index().
		Index(indexName).
		Id(docID).
		BodyJson(document).
		Do(ctx)
	if err != nil {
		s.logger.Error("Failed to index document %s in index %s: %v", docID, indexName, err)
		return fmt.Errorf("failed to index document: %w", err)
	}
	s.logger.Info("Indexed document %s in index %s", docID, indexName)
	return nil
}

func (s *genericSearchStore) Search(ctx context.Context, indexName string, query string, fields []string, page, pageSize int) ([]json.RawMessage, int64, error) {
	from := (page - 1) * pageSize
	if from < 0 {
		from = 0
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	multiMatchQuery := elastic.NewMultiMatchQuery(query, fields...).
		Type("best_fields").
		Fuzziness("AUTO").
		MinimumShouldMatch("75%")

	searchResult, err := s.esClient.ESClient.Search().
		Index(indexName).
		Query(multiMatchQuery).
		From(from).
		Size(pageSize).
		Do(ctx)

	if err != nil {
		s.logger.Error("Failed to search in index %s: %v", indexName, err)
		return nil, 0, fmt.Errorf("failed to search: %w", err)
	}

	var rawMessages []json.RawMessage
	if searchResult.Hits != nil && searchResult.Hits.Hits != nil {
		for _, hit := range searchResult.Hits.Hits {
			rawMessages = append(rawMessages, hit.Source)
		}
	}

	s.logger.Info("Found %d documents for query '%s' in index '%s'", searchResult.Hits.TotalHits.Value, query, indexName)
	return rawMessages, searchResult.Hits.TotalHits.Value, nil
}

func (s *genericSearchStore) DeleteDocument(ctx context.Context, indexName string, docID string) error {
	_, err := s.esClient.ESClient.Delete().
		Index(indexName).
		Id(docID).
		Do(ctx)
	if err != nil {
		s.logger.Error("Failed to delete document %s from index %s: %v", docID, indexName, err)
		return fmt.Errorf("failed to delete document: %w", err)
	}
	s.logger.Info("Deleted document %s from index %s", docID, indexName)
	return nil
}
