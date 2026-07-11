package config

import (
	"os"
	"runtime"
	"strconv"
)

type Config struct {
	WorkerCount  int
	JobQueueSize int
	MetricsPort  string
}

func Load() Config {

	return Config{
		WorkerCount:  getIntEnv("WORKER_COUNT", runtime.NumCPU()),
		JobQueueSize: getIntEnv("JOB_QUEUE_SIZE", 100),
		MetricsPort:  getStringEnv("METRICS_PORT", "2112"),
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

	return intValue
}

func getStringEnv(key string, defaultValue string) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
