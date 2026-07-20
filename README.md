# Metadata Ingestion Service


![Go](https://img.shields.io/badge/Go-1.24-blue)

![License](https://img.shields.io/badge/license-MIT-green)

![Docker](https://img.shields.io/badge/docker-ready-blue)

![Prometheus](https://img.shields.io/badge/prometheus-enabled-orange)

Metadata Ingestion Service is a cloud-native backend application written in Go that demonstrates production-grade software engineering practices including concurrent processing, worker pools, retry mechanisms, observability, graceful shutdown, and scalable connector-based ingestion pipelines.

Inspired by enterprise metadata platforms, the project showcases how to build maintainable, fault-tolerant backend services using modern Go design principles.

## Why this project?

Enterprise metadata platforms ingest millions of metadata entities from external systems such as Power BI, Tableau, MLflow, databases, and cloud services.

This project demonstrates how such ingestion pipelines can be designed using scalable Go concurrency patterns while maintaining reliability, observability, and clean architecture.

Although simplified for demonstration purposes, the architecture closely follows patterns used in production systems.

The goal of this project is to demonstrate production-ready backend engineering concepts including:

- Concurrent metadata ingestion
- Connector framework
- Worker pools
- Context propagation
- Graceful shutdown
- Retry mechanisms
- Rate limiting
- Observability
- Profiling

---

## Architecture

                    +----------------------+
                    |      HTTP Server     |
                    +----------+-----------+
                               |
                     Middleware (Logging, Recovery)
                               |
                               ▼
                     Ingestion Service
                               |
                        Worker Pool
                 (goroutines + channels)
                               |
            +------------------+------------------+
            |                  |                  |
         Power BI          Tableau           MLflow
            |                  |                  |
            +------------------+------------------+
                               |
                           Processor
                               |
                       Retry Framework
                               |
                            Sink Layer
                               |
                           Console(OpenSearch)
                               |
              Prometheus Metrics + OpenTelemetry

---

## Tech Stack

### Language

- Go 1.26

### Libraries

- Zap
- Prometheus Client
- OpenTelemetry

### Infrastructure

- Docker

### CI/CD

- GitHub Actions

### Design Patterns

- Worker Pool
- Dependency Injection
- Strategy Pattern
- Retry Pattern

---

## Project Structure

```text
cmd/
internal/
    connectors/
    processor/
    worker/
    retry/
    metrics/
    middleware/
pkg/
configs/
docker/
```

---

## API Documentation

The Metadata Ingestion Service exposes a minimal REST API for monitoring and triggering metadata ingestion workflows.

**Base URL**

```text
http://localhost:8080
```

---

## Endpoints

| Method | Endpoint | Description |
|---------|----------|-------------|
| GET | `/health` | Returns the liveness status of the service |
| GET | `/ready` | Returns the readiness status of the service |
| POST | `/ingest` | Starts an asynchronous metadata ingestion job |

---

## GET /health

Returns the current health status of the application.

### Request

```http
GET /health HTTP/1.1
Host: localhost:8080
```

### Response

**Status:** `200 OK`

```json
{
  "status": "UP"
}
```

### Description

This endpoint is intended for:

- Kubernetes Liveness Probe
- Docker health checks
- Load balancer monitoring
- Service availability verification

---

## GET /ready

Returns whether the service is ready to accept ingestion requests.

### Request

```http
GET /ready HTTP/1.1
Host: localhost:8080
```

### Response

**Status:** `200 OK`

```json
{
    "status": "READY"
}
```

### Description

This endpoint is intended for:

- Kubernetes Readiness Probe
- Deployment rollouts
- Traffic routing
- Startup verification

---

## POST /ingest

Starts a metadata ingestion workflow.

The ingestion process executes asynchronously. The API immediately returns an acknowledgment while the ingestion pipeline continues processing in the background.

### Request

```http
POST /ingest HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

No request body is required.

### Response

**Status:** `202 Accepted`

```json
{
    "message": "Ingestion started"
}
```

### Processing Workflow

Once an ingestion request is received, the service performs the following steps:

1. Starts the ingestion service in a background goroutine.
2. Creates a buffered job queue.
3. Initializes a configurable worker pool.
4. Retrieves metadata from all configured connectors.
5. Converts metadata into ingestion jobs.
6. Places jobs into the queue.
7. Workers process jobs concurrently.
8. Processed metadata is written to the configured sink.
9. Prometheus metrics and application logs are generated throughout the execution.
10. The service waits for all workers to complete before terminating the ingestion run.

---

## Request Lifecycle

```text
Client
   │
   ▼
POST /ingest
   │
   ▼
HTTP Handler
   │
   ▼
Start Background Goroutine
   │
   ▼
Create Job Queue
   │
   ▼
Start Worker Pool
   │
   ▼
Fetch Metadata
(Power BI, Tableau, MLflow)
   │
   ▼
Create Metadata Jobs
   │
   ▼
Buffered Channel
   │
   ▼
Concurrent Workers
   │
   ▼
Processor
   │
   ▼
Sink (OpenSearch)
```

---

## Concurrency Model

The ingestion service uses a worker pool pattern for efficient concurrent processing.

- Configurable worker count
- Buffered job queue
- Goroutines for parallel execution
- Channels for job distribution
- WaitGroup for synchronization
- Context propagation for graceful shutdown
- Automatic cleanup after processing completes

---

## HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| `200 OK` | Health or readiness check completed successfully |
| `202 Accepted` | Metadata ingestion has been started successfully |
| `500 Internal Server Error` | Unexpected server error while processing the request |

---

## Notes

- Metadata ingestion is asynchronous and non-blocking.
- Worker count and job queue size are configurable through the application configuration.
- All configured connectors are processed during an ingestion run.
- Context propagation enables graceful shutdown and cancellation of in-flight operations.
- The service follows a connector-based architecture, making it easy to add support for additional metadata sources.
- Prometheus metrics and structured logging provide operational observability.

---

## Features

- Connector abstraction
- Concurrent metadata ingestion
- Worker Pool implementation
- Context propagation
- Graceful shutdown
- Retry with exponential backoff
- Dependency Injection
- Prometheus metrics
- HTTP health endpoints
- Docker support
- GitHub Actions CI

## Concepts Demonstrated

- Interfaces
- Dependency Injection
- Goroutines
- Channels
- WaitGroups
- Context
- Worker Pools
- Middleware
- Retry Pattern
- Prometheus Metrics
- HTTP Servers
- Docker
- CI/CD

## Lessons Learned

Building this project reinforced several key backend engineering principles:

- Concurrency should improve throughput without sacrificing maintainability.
- Context propagation is essential for graceful cancellation.
- Retry logic must be carefully designed to avoid overwhelming downstream services.
- Observability should be built in from the beginning rather than added later.
- Interfaces make connector implementations easier to extend and test.