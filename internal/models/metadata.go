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
