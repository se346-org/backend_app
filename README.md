# Chat Service Backend

A backend service for a real-time chat application built with Go.

## Features
- [x] User registration and login
- [x] JWT-based authentication  
- [x] WebSocket support for real-time updates
- [x] Direct Messages (DM)
- [x] Structured logging with request tracing
- [x] Prometheus metrics collection
- [x] Distributed tracing with OpenTelemetry
- [x] Health checks and monitoring
- [ ] Group Chat
- [ ] Reply to messages

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 16.x
- Redis (for real-time features)
- Nats (for pubsub)
- Prometheus (for metrics)
- Jaeger (for tracing)

## Configuration

Create a `config.yaml` file in the root directory with the following variables:

```yaml
server:
    port: 8080
    origin: "http://localhost:8080"
postgres:
    host: localhost
    port: 5432
    username: postgres
    password: password
    database: chat_socio
    sslmode: disable
redis:
    host: localhost
    port: 6379
    password: ""
    database: 0
nats:
    address: "nats://localhost:4222"
jwt:
    secret: "ASDWRTGHJKLQWERTYUIOPZXCVBNMASDFGHJKLQWERTYUIOPZXCVBNM"
    issuer: "chat_socio"
    expiration: 3600
observability:
    tracing_enabled: true
    jaeger_endpoint: localhost:4318
    jaeger_service: "chat-socio"
```

## Run

### Using Docker Compose (Recommended)
```bash
# Start all services including observability stack
docker-compose up -d

# View logs
docker-compose logs -f chat-service
```

### Manual Setup
```bash
# Start infrastructure services
docker-compose up -d postgres redis nats prometheus jaeger grafana

# Run the application
go run cmd/app/main.go
```

## Observability

### Metrics
- Application metrics are exposed at `http://localhost:8080/metrics`
- Prometheus dashboard: `http://localhost:9090`
- Grafana dashboard: `http://localhost:3000` (admin/admin)

### Tracing
- Jaeger UI: `http://localhost:16686`
- Traces include HTTP requests, database queries, and message processing

### Logging
- Structured JSON logging with correlation IDs
- Log levels: DEBUG, INFO, WARN, ERROR
- Request/response logging with timing


## Technology Stack
- **Runtime**: Go 1.24
- **Database**: PostgreSQL 16.x
- **Cache**: Redis 7.x
- **Message Queue**: Nats 2.9.x
- **HTTP Server**: Hertz (https://github.com/cloudwego/hertz)
- **Authentication**: JWT
- **Observability**: OpenTelemetry with Jaeger backend

