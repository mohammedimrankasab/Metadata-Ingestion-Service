package models

import "time"

type MetadataType string

const (
	ReportType    MetadataType = "REPORT"
	DashboardType MetadataType = "DASHBOARD"
	DatasetType   MetadataType = "DATASET"
)

type Metadata struct {
	ID           string
	Name         string
	Type         MetadataType
	Workspace    string
	Source       string
	LastModified time.Time
}

func NewMetadata(
	id string,
	name string,
	metadataType MetadataType,
	workspace string,
	source string,
	lastModified time.Time,
) Metadata {
	return Metadata{
		ID:           id,
		Name:         name,
		Type:         metadataType,
		Workspace:    workspace,
		Source:       source,
		LastModified: lastModified,
	}
}
