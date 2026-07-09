# Metadata Ingestion Service

A production-style metadata ingestion framework written in Go, inspired by enterprise metadata platforms.

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
- Kubernetes readiness

---

## Architecture

> 🚧 Architecture diagram will be added as the project evolves.

---

## Roadmap

| Sprint | Status | Description |
|---------|--------|-------------|
| Sprint 1 | ✅ Completed | Connector framework, metadata model, Power BI connector |
| Sprint 2 | ✅ Completed | Concurrent ingestion engine using worker pools |
| Sprint 3 | ⏳ Next | Retry, backoff and graceful shutdown |
| Sprint 4 | ⏳ Planned | Metrics, profiling and observability |
| Sprint 5 | ⏳ Planned | Scheduler and background jobs |

---

## Milestones

### ✅ Milestone 1 - Connector Framework

Implemented:

- Connector abstraction
- Power BI connector
- Metadata model
- Ingestion service
- Context support
- Incremental sync support (`lastSyncTime`)

Git Commit:

```text
git commit -m "refactor: introduce application composition root"
```

---

## Tech Stack

- Go
- Cobra
- Viper
- Zap
- Docker
- GitHub Actions (Upcoming)

---

## Upcoming Features

- Worker Pool
- Buffered Channels
- Context Cancellation
- Graceful Shutdown
- OpenSearch Sink
- REST API
- Metrics
- Profiling

### V1 Architecture proposal

```
                         Metadata Ingestion Service

                   +-------------------------------+
                   |            cmd/app            |
                   +---------------+---------------+
                                   |
                                   |
                         Ingestion Service
                                   |
                  +----------------+----------------+
                  |                                 |
           Connector Manager                 Worker Pool
                  |                                 |
        +---------+---------+                Metadata Channel
        |         |         |                       |
     PowerBI   Tableau    MLflow                    |
                  |                                 |
                  +---------------+-----------------+
                                  |
                           Metadata Processor
                                  |
                           OpenSearch Sink
```