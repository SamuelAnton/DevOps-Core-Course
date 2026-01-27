# Overview
A Go implementation of the DevOps Info Service with the same functionality as the Python/Flask version.

# Prerequisites
- Go 1.21 or higher

# Installation
1. **Clone the repository**:
2. **Build**:
```bash
# Build for current platform
go build -o devops-info-service main.go

# Build for Linux (cross-compilation)
GOOS=linux GOARCH=amd64 go build -o devops-info-service-linux main.go

```

# Running the Application
```
# Run directly with go
go run main.go

# Run compiled binary
./devops-info-service

# Run with custom configuration
HOST=127.0.0.1 PORT=3000 DEBUG=true ./devops-info-service
```

# API Endpoint
## `GET /` - Service and system information
### Example Response:
```
{
  "service": {
    "name": "devops-info-service",
    "version": "1.0.0",
    "description": "DevOps course info service",
    "framework": "Go"
  },
  "system": {
    "hostname": "my-laptop",
    "platform": "Linux",
    "platform_version": "Ubuntu 24.04",
    "architecture": "x86_64",
    "cpu_count": 8,
    "go": "go1.21.0"
  },
  "runtime": {
    "uptime_seconds": 3600,
    "uptime_human": "1 hour, 0 minutes",
    "current_time": "2026-01-07T14:30:00.000Z",
    "timezone": "UTC"
  },
  "request": {
    "client_ip": "127.0.0.1",
    "user_agent": "curl/7.81.0",
    "method": "GET",
    "path": "/"
  },
  "endpoints": [
    {"path": "/", "method": "GET", "description": "Service information"},
    {"path": "/health", "method": "GET", "description": "Health check"}
  ]
}
```

## `GET /health` - Health check
### Example Response:
```
{
  "status": "healthy",
  "timestamp": "2024-01-15T14:30:00.000Z",
  "uptime_seconds": 3600
}
```

# Configuration 
| Variable name | Basic value | Description |
|----------|--------|-------------|
| **Host** | 0.0.0.0 | IP of the service |
| **Port** | 5000 | Port of the service |
| **Debug** | False | Debug mode enabeled |
