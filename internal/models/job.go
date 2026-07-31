package models

import (
	"time"

	"github.com/google/uuid"
)

type MetadataJob struct {
	ID        string
	Metadata  Metadata
	Connector string
	CreatedAt time.Time
	Done      chan struct{}
}

func NewJob(
	connector string,
	metadata Metadata,
) MetadataJob {
	return MetadataJob{
		ID:        uuid.NewString(),
		Metadata:  metadata,
		Connector: connector,
		CreatedAt: time.Now(),
	}
}
