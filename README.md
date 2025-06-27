# Teleprobe - Telemetry Sensor

## Overview

**Teleprobe** is a telemetry sensor system designed to collect, process, and transmit telemetry data efficiently. This project implements sensor nodes and a telemetry sink, facilitating the ingestion, aggregation, and analysis of telemetry data streams.

The sensor nodes gather data from various sources and send it to the telemetry sink for storage and further processing.

## Architecture Overview

```text
+------------------+     gRPC / HTTP / MQTT / NATS     +--------------------+
|  Sensor Node     | --------------------------------> |  Telemetry Sink    |
|  (Data Emitter)  |                                   |  (Data Receiver)   |
+------------------+                                   +--------------------+
        ↑                                                      ↓
        |                                              DB / File / Message Queue
        | (Retry, Buffer)                                   (Storage, Analytics)
        |
    Local Buffer
```

- **Sensor Nodes**: Collect telemetry data and transmit it using protocols like gRPC or HTTP.
- **Telemetry Sink**: Receives data streams, aggregates metrics, and stores or forwards the data.

## Features

- Modular sensor nodes supporting multiple telemetry sources.
- Dockerized components for easy deployment.

## Getting Started

### Prerequisites

- Go 1.20+ installed
- Docker and Docker Compose (optional, for containerized setup)
- Make (optional, for build automation)

### Installation

Clone the repository:

```bash
git clone https://github.com/ximura/teleprobe.git
cd teleprobe
```

## Run Instructions

### ✅ Local Development

1. **Set environment variables** (optionally via `.env` file):

```env
BIND_ADDR=:50051
LOG_FILE=telemetry.log
BUFFER_SIZE=1000
```

2. **Start the sink service**:

```bash
go run ./cmd/sink
```

3. **Start a sensor instance**:

```bash
go run ./cmd/sensor
```

Each service has its own configuration loaded from environment variables.

### 🐳 Docker Deployment

1. **Build Docker images**:

```bash
docker build -f Dockerfile.sensor -t teleprobe-sensor .
docker build -f Dockerfile.sink -t teleprobe-sink .
```

2. **Start via Docker Compose**:

```bash
docker-compose up --build
```

## Testing

Run unit tests:

```bash
go test ./...
```

## Improvements

### Sensor Enhancements

- **Dropped Message Caching**: Currently, if the sink is unavailable, telemetry messages may be lost. Implementing a local caching mechanism (e.g., in-memory buffer, local file queue) would allow the sensor to retry sending messages when the connection is restored.
- **Retry Logic**: Introducing an exponential backoff retry strategy for failed transmissions would make the sensor more resilient to temporary network issues.

### Communication Layer

- **Message Bus Integration**: Replacing the direct gRPC connection with a message broker such as **RabbitMQ** or **Kafka** would offer:
  - Built-in message persistence
  - Retry support
  - Dead Letter Queue (DLQ) mechanisms for debugging undeliverable telemetry
  - Better decoupling between sensor and sink services

### Sink Storage Strategy

- **File Storage Limitation**: File-based telemetry storage is difficult to scale and query.
- **Timeseries Database Recommendation**: For scalability and efficient querying, switching to a timeseries database like **VictoriaMetrics**, **InfluxDB** would be a significant improvement.

These improvements would make the system more reliable, maintainable, and production-ready.