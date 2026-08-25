# Balances Worker

> Go worker responsible for processing balance events and generating daily balance reports.

## Overview

Balances Worker is a Go-based background service responsible for processing balance-related events and generating the daily balance report.

The project complements the Balances API by moving report processing away from the synchronous application flow and into a dedicated worker.

This separation allows the processing workload to evolve independently from the API responsible for balance operations.

## Architecture

```text
                  ┌───────────────────┐
                  │    Balances API   │
                  └─────────┬─────────┘
                            │
                            │ Balance Events
                            ▼
                  ┌───────────────────┐
                  │   Event Source    │
                  └─────────┬─────────┘
                            │
                            │
                            ▼
                  ┌───────────────────┐
                  │  Balances Worker  │
                  │                   │
                  │        Go         │
                  └─────────┬─────────┘
                            │
                            │ Process
                            ▼
                  ┌───────────────────┐
                  │ Daily Balance     │
                  │     Report        │
                  └───────────────────┘
```

The worker is intentionally isolated from the API layer so that balance-report processing can run independently.

## Responsibilities

The main responsibility of the worker is to process balance events and generate the daily balance report.

The service is responsible for:

- Processing balance-related events
- Executing background processing
- Generating the daily balance report
- Isolating report processing from the API lifecycle

## Why a Worker?

Generating reports can involve processing multiple events and records and does not necessarily need to happen during an HTTP request.

Moving this workload to a background worker provides a clear separation between:

```text
Synchronous API Operations
            │
            ▼
       Balances API
            │
            │
            ▼
   Asynchronous Processing
            │
            ▼
     Balances Worker
            │
            ▼
    Daily Balance Report
```

This architecture allows the API and reporting workload to scale independently.

## Project Structure

```text
balances-worker/
│
├── config/
│   └── Application configuration
│
├── internal/
│   └── worker/
│       └── Worker implementation
│
├── scripts/
│   └── docker/
│       └── Docker resources
│
├── .gitignore
├── go.mod
└── README.md
```

### `config/`

Contains configuration required by the worker.

### `internal/worker/`

Contains the worker implementation and processing logic.

Keeping the worker implementation under `internal` prevents it from being imported directly by external packages.

### `scripts/docker/`

Contains Docker-related resources used to build or run the worker.

## Processing Model

The worker follows a background-processing model:

```text
            Event
              │
              ▼
      ┌───────────────┐
      │    Worker     │
      └───────┬───────┘
              │
              ▼
      Process Balance
          Events
              │
              ▼
       Aggregate Data
              │
              ▼
      Generate Daily
          Report
```

This model separates event ingestion from report generation and allows processing to happen independently of the API request lifecycle.

## Technology Stack

| Technology | Purpose |
|---|---|
| Go | Worker implementation |
| Docker | Containerization |

## Engineering Concepts

This project demonstrates practical backend engineering concepts including:

- Go backend development
- Background processing
- Worker-based architecture
- Asynchronous processing
- Event-driven architecture
- Separation of concerns
- Service isolation
- Containerization

## Relationship with Balances API

Balances Worker is designed as a companion service to the Balances API.

```text
                  ┌───────────────────┐
                  │    Balances API   │
                  │                   │
                  │ Balance Operations│
                  └─────────┬─────────┘
                            │
                         Events
                            │
                            ▼
                  ┌───────────────────┐
                  │  Balances Worker  │
                  │                   │
                  │ Event Processing   │
                  │       +           │
                  │ Report Generation  │
                  └───────────────────┘
```

The separation allows the API to focus on balance operations while the worker handles background reporting workloads.

## Local Development

### Requirements

- Go
- Docker

### Clone

```bash
git clone https://github.com/clodoaldomarques/balances-worker.git
cd balances-worker
```

### Install dependencies

```bash
go mod download
```

### Run tests

```bash
go test ./...
```

### Build

```bash
go build ./...
```

### Run

```bash
go run ./...
```

## Docker

Docker resources are available under:

```text
scripts/docker/
```

Build the application using the project's Docker configuration.

## Design Principles

The project follows several principles commonly used in distributed backend systems:

- **Single responsibility** — report processing is isolated in a dedicated service.
- **Asynchronous execution** — background workloads do not need to block API requests.
- **Service isolation** — reporting concerns are separated from balance operations.
- **Scalability** — worker instances can potentially be scaled independently from the API.
- **Maintainability** — worker-specific logic remains isolated from other services.

## Project Status

This project is part of my Go backend engineering portfolio and demonstrates the implementation of a dedicated worker for asynchronous balance-event processing and daily reporting.

The project is intended primarily as an engineering study and portfolio project.

## Author

**Clodoaldo Marques**

Backend Software Engineer focused on Go, Microservices, Distributed Systems and Cloud-Native architectures.

- GitHub: https://github.com/clodoaldomarques
- LinkedIn: https://www.linkedin.com/in/clodoaldomarques/