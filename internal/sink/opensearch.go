package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

type OpenSearchSink struct {
	logger *zap.Logger
	client *http.Client
	url    string
	index  string
}

func NewOpenSearchSink(
	logger *zap.Logger,
	url string,
	index string,
) *OpenSearchSink {

	return &OpenSearchSink{
		logger: logger,
		client: &http.Client{},
		url:    url,
		index:  index,
	}
}

func (s *OpenSearchSink) Write(
	ctx context.Context,
	metadata models.Metadata,
) error {

	body, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s/_doc", s.url, s.index),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("opensearch returned status %d", resp.StatusCode)
	}

	s.logger.Debug(
		"metadata indexed",
		zap.String("name", metadata.Name),
	)

	return nil
}
