# Concurrent Job Processor

A concurrent background job processing system built with Go that demonstrates Go's concurrency primitives through a worker pool architecture.

## Features

- REST API for submitting and tracking background jobs
- Concurrent job processing using goroutines
- Buffered channels for job queuing
- Worker pool architecture
- Job status tracking
- Job cancellation using Go context
- Job processing timeouts
- Thread-safe context management using mutexes
- Graceful shutdown using WaitGroups
- MySQL database persistence
- Docker and Docker Compose support

## Supported Job Types

- `generate_report`
- `send_email`
- `process_data`

## Job Lifecycle

Jobs can move through the following states:

`queued → processing → completed`

Jobs may also end with:

`failed`

or:

`cancelled`

## Tech Stack

- Go
- Gin
- GORM
- MySQL
- Goroutines
- Channels
- Context
- Sync Mutex
- Sync WaitGroup
- Docker
- Docker Compose

## API Endpoints

### Health Check

`GET /`

### Create Job

`POST /api/v1/jobs`

Example request:

```json
{
  "type": "generate_report",
  "payload": "Generate monthly sales report"
}