# Metadata Ingestion Service


![Go](https://img.shields.io/badge/Go-1.24-blue)

![License](https://img.shields.io/badge/license-MIT-green)

![Docker](https://img.shields.io/badge/docker-ready-blue)

![Kubernetes](https://img.shields.io/badge/kubernetes-ready-blue)

![Prometheus](https://img.shields.io/badge/prometheus-enabled-orange)

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
                          OpenSearch
                               |
              Prometheus Metrics + OpenTelemetry

---

## Tech Stack

- Go
- Zap
- Docker
- GitHub Actions

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
- Kubernetes manifests
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
- Kubernetes
- CI/CD