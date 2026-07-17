# Concurrent Job Processor

A concurrent background job processing system built with **Go**, designed to demonstrate practical usage of Go's concurrency primitives including **goroutines, channels, worker pools, mutexes, contexts, and WaitGroups**.

The system allows clients to submit background jobs through REST APIs. Jobs are placed into a buffered channel and processed concurrently by a fixed-size pool of worker goroutines. Job states are persisted in MySQL and can be tracked, cancelled, or queried through the API.

---

## Features

- RESTful API built with Gin
- Asynchronous background job processing
- Concurrent processing using Goroutines
- Buffered Channels for job queuing
- Fixed-size Worker Pool architecture
- MySQL database persistence using GORM
- Job status tracking
- Job failure handling
- Job cancellation using Go Context
- Job processing timeouts
- Thread-safe context management using Mutex
- Graceful shutdown using WaitGroup
- Dockerized backend and MySQL database
- Docker Compose for multi-container setup
- Layered backend architecture

---

## Architecture

```text
Client
   │
   │ POST /api/v1/jobs
   ▼
Gin REST API
   │
   ▼
Service Layer
   │
   ├──────────────► MySQL Database
   │                 status = queued
   │
   ▼
Buffered Job Channel
   │
   ├────────────┬────────────┐
   ▼            ▼            ▼
Worker 1     Worker 2     Worker 3
Goroutine    Goroutine    Goroutine
   │            │            │
   ▼            ▼            ▼
Process Job  Process Job  Process Job
   │            │            │
   └────────────┴────────────┘
                │
                ▼
          MySQL Database
                │
       ┌────────┼────────┐
       ▼        ▼        ▼
  completed   failed   cancelled
```

When a client submits a job, the API immediately creates a database record with the `queued` status and sends the job to a buffered Go channel.

Multiple worker goroutines listen to this channel and process jobs concurrently. Each worker updates the job status as processing progresses.

---

## Job Lifecycle

A successfully processed job follows:

```text
queued → processing → completed
```

A job that encounters an error follows:

```text
queued → processing → failed
```

A running or queued job can also be cancelled:

```text
queued/processing → cancelled
```

---

## Concurrency Concepts Demonstrated

### Goroutines

Each worker runs as an independent goroutine, allowing multiple background jobs to be processed concurrently.

### Channels

A buffered Go channel acts as an in-memory job queue that distributes submitted jobs among available workers.

### Worker Pool

The application starts a fixed number of workers instead of creating an unlimited number of goroutines.

This provides controlled concurrency and prevents uncontrolled resource usage.

### Mutex

A mutex protects the shared map containing job cancellation functions, ensuring safe concurrent access from multiple goroutines.

### Context

Go's `context` package is used to support job cancellation and processing timeouts.

A running job can stop processing when its context is cancelled.

### WaitGroup

A `sync.WaitGroup` tracks active worker goroutines during application shutdown.

The application waits for workers to finish before terminating.

### Graceful Shutdown

When the application receives an interrupt or termination signal:

1. The HTTP server stops accepting new requests.
2. The job queue is closed.
3. Workers finish jobs already being processed.
4. The application waits for all workers to exit.
5. The application shuts down safely.

---

## Supported Job Types

The system currently supports the following simulated background jobs:

| Job Type | Description | Processing Time |
|---|---|---|
| `send_email` | Simulates sending an email | ~3 seconds |
| `generate_report` | Simulates report generation | ~6 seconds |
| `process_data` | Simulates data processing | ~8 seconds |

Unsupported job types are automatically marked as `failed`.

---

## Tech Stack

| Technology | Purpose |
|---|---|
| Go | Backend programming language |
| Gin | HTTP web framework |
| GORM | ORM for database operations |
| MySQL | Persistent job storage |
| Goroutines | Concurrent worker execution |
| Channels | In-memory job queue |
| Context | Cancellation and timeout handling |
| Mutex | Thread-safe shared state |
| WaitGroup | Worker synchronization |
| Docker | Application containerization |
| Docker Compose | Multi-container orchestration |

---

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/` | Health check |
| POST | `/api/v1/jobs` | Submit a new background job |
| GET | `/api/v1/jobs` | Get all jobs |
| GET | `/api/v1/jobs/:id` | Get a job by ID |
| DELETE | `/api/v1/jobs/:id` | Cancel a job |

---

## API Examples

### Health Check

```http
GET /
```

Example response:

```json
{
  "message": "Concurrent Job Processor API is running",
  "success": true
}
```

---

### Submit a Job

```http
POST /api/v1/jobs
```

Request body:

```json
{
  "type": "generate_report",
  "payload": "Generate monthly sales report"
}
```

Example response:

```json
{
  "success": true,
  "message": "Job submitted successfully",
  "data": {
    "id": 1,
    "type": "generate_report",
    "payload": "Generate monthly sales report",
    "status": "queued"
  }
}
```

The job is then picked up asynchronously by an available worker.

---

### Get Job by ID

```http
GET /api/v1/jobs/1
```

A completed job may return:

```json
{
  "success": true,
  "message": "Job fetched successfully",
  "data": {
    "id": 1,
    "type": "generate_report",
    "payload": "Generate monthly sales report",
    "status": "completed",
    "result": "Report generated successfully: Generate monthly sales report"
  }
}
```

---

### Get All Jobs

```http
GET /api/v1/jobs
```

Returns all submitted jobs and their current statuses.

---

### Cancel a Job

```http
DELETE /api/v1/jobs/1
```

Running jobs are cancelled using Go's Context cancellation mechanism.

Queued jobs are marked as cancelled and skipped when picked up by a worker.

---

## Project Structure

```text
concurrent-job-processor/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── config/
│   └── database.go
│
├── internal/
│   ├── handlers/
│   │   └── job_handler.go
│   │
│   ├── models/
│   │   ├── job.go
│   │   └── job_request.go
│   │
│   ├── repository/
│   │   └── job_repository.go
│   │
│   ├── routes/
│   │   └── job_routes.go
│   │
│   ├── services/
│   │   └── job_service.go
│   │
│   └── worker/
│       ├── context_manager.go
│       ├── pool.go
│       ├── queue.go
│       └── worker.go
│
├── pkg/
│   └── response/
│       └── response.go
│
├── .dockerignore
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## Running Locally

### Prerequisites

Make sure you have installed:

- Go
- MySQL

Create the database:

```sql
CREATE DATABASE job_processor;
```

Create a `.env` file based on `.env.example` and configure your database credentials.

Example:

```env
SERVER_PORT=8081

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=job_processor
```

Run the application:

```bash
go run ./cmd/server
```

The API will be available at:

```text
http://localhost:8081
```

---

## Running with Docker

The application can also be run using Docker Compose.

Docker Compose starts:

- Go backend container
- MySQL database container

Build and start the containers:

```bash
docker compose up --build
```

The API will be available on port `8081`.

Stop the containers:

```bash
docker compose down
```

MySQL data is persisted using a Docker volume.

---

## Environment Variables

| Variable | Description |
|---|---|
| `SERVER_PORT` | Port used by the HTTP server |
| `DB_HOST` | MySQL database host |
| `DB_PORT` | MySQL database port |
| `DB_USER` | MySQL username |
| `DB_PASSWORD` | MySQL password |
| `DB_NAME` | MySQL database name |

Use `.env.example` as a reference.

The actual `.env` file is excluded from Git using `.gitignore`.

---

## Future Improvements

Possible improvements for a production-grade implementation include:

- Redis or RabbitMQ for a durable job queue
- Automatic retry mechanism for failed jobs
- Exponential backoff
- Job priorities
- Dynamic worker scaling
- Authentication and authorization
- Prometheus metrics and monitoring
- Structured logging
- Distributed worker support
- Dead-letter queues for permanently failed jobs

---

## Key Learning Outcomes

This project demonstrates how Go can be used to build concurrent backend systems using:

- Goroutines for lightweight concurrent execution
- Channels for communication between components
- Worker pools for controlled concurrency
- Context for cancellation and timeouts
- Mutexes for synchronization
- WaitGroups for coordinating goroutines
- Graceful shutdown for safe application termination

It also demonstrates REST API development, layered backend architecture, MySQL persistence, and Docker-based deployment.