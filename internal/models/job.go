package models

import "github.com/google/uuid"

type MetadataJob struct {
	ID       string
	Metadata Metadata
}

func NewJob(metadata Metadata) MetadataJob {
	return MetadataJob{
		ID:       uuid.NewString(),
		Metadata: metadata,
	}
}
