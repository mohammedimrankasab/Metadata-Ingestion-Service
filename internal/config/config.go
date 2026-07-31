// Package config provides application configuration loaded from
// environment variables with sensible defaults.
package config

import (
	"os"
	"runtime"
	"strconv"
)

type Config struct {
	WorkerCount  int
	JobQueueSize int
	HTTPPort     string

	SinkType   string
	OutputFile string

	OpenSearchURL   string
	OpenSearchIndex string
	EnableTracing   bool
}

func Load() *Config {

	return &Config{
		WorkerCount:  getIntEnv("WORKER_COUNT", runtime.NumCPU()),
		JobQueueSize: getIntEnv("JOB_QUEUE_SIZE", 100),
		HTTPPort:     getStringEnv("HTTP_PORT", "8080"),

		SinkType:   getStringEnv("SINK_TYPE", "console"),
		OutputFile: getStringEnv("OUTPUT_FILE", "metadata.json"),

		OpenSearchURL:   getStringEnv("OPENSEARCH_URL", "http://localhost:9200"),
		OpenSearchIndex: getStringEnv("OPENSEARCH_INDEX", "metadata"),

		EnableTracing: getBoolEnv("ENABLE_TRACING", false),
	}
}
func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	if intValue <= 0 {
		return defaultValue
	}

	return intValue
}

func getStringEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getBoolEnv(key string, defaultValue bool) bool {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
