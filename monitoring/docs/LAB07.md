# LAB 07

## Architecture
The logging stack consists of the following components:
- Loki – log aggregation system (stores and indexes logs)
- Promtail – log collector (scrapes container logs and pushes to Loki)
- Grafana – visualization and query interface
- Application containers (Python and Go) – generate structured JSON logs

All components run as Docker containers managed by Docker Compose. They communicate over a dedicated Docker network (logging).

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│   Application   │      │    Promtail     │      │      Loki       │
│  (Python/Go)    │─────▶│  (log scraper)  │─────▶│ (log storage)   │
└─────────────────┘      └─────────────────┘      └─────────────────┘
         │                                                    │
         │ (logs via Docker)                                  │ (API)
         ▼                                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                              Grafana                                │
│                      (query and visualization)                      │
└─────────────────────────────────────────────────────────────────────┘
```

Flow:
1. Applications write logs to stdout/stderr in JSON format.
2. Docker captures these logs and stores them in JSON files on the host.
3. Promtail (running with access to Docker socket and log directories) discovers containers via Docker
4. Service Discovery and tails their log files.
5. Promtail sends log entries to Loki’s HTTP API.

Grafana connects to Loki as a data source and allows querying logs using LogQL.

## Setup Guide
### Prerequisites
- Docker and Docker Compose installed
- Git (optional)
- Clone of the project repository containing the monitoring/ directory

### Step 1: Prepare Configuration Files
Create the following directory structure and files as described in previous tasks. Ensure you have:
- monitoring/docker-compose.yml
- monitoring/loki/config.yml
- monitoring/promtail/config.yml
- monitoring/.env (for Grafana admin password)

### Step 2: Start the Stack
```
cd monitoring
docker compose up -d
docker compose ps
```
All containers should be running.

### Step 3: Verify Component Health
```
curl http://localhost:3100/ready          # Loki ready
curl http://localhost:9080/targets        # Promtail targets
curl http://localhost:3000                # Grafana (login page)
```

### Step 4: Configure Grafana Data Source
1. Log in to Grafana at http://localhost:3000 (admin / password from .env).
2. Go to Connections → Data sources → Add data source.
3. Choose Loki.
4. Set URL to http://loki:3100 (container name, because Grafana is on same network).
5. Click Save & Test. You should see a success message.

### Step 5: Generate Test Logs
```
for i in {1..20}; do curl http://localhost:8000/; done
for i in {1..20}; do curl http://localhost:8000/health; done
for i in {1..20}; do curl http://localhost:8001/; done
```

### Step 6: Explore Logs in Grafana
- Go to Explore.
- Select Loki data source.
- Use the Label browser or type queries like {app="devops-python"}.
- Verify that logs appear.

## Configuration
### Loki Configuration (loki/config.yml)
```
schema_config:
  configs:
    - from: 2024-01-01
      store: tsdb        # 🚀 New fast store
      object_store: filesystem
      schema: v13        # 📊 Latest schema
      index:
        prefix: index_
        period: 24h

auth_enabled: false

server:
  http_listen_port: 3100

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory

limits_config:
  retention_period: 168h  # 🗓️ 7 days
```

Why this configuration?
- `auth_enabled`: false – no authentication for internal use (simplifies setup).
- `filesystem` storage` – suitable for single-node testing.
- `tsdb` with schema v13 – modern, efficient index.
- `retention_period: 168h` – logs kept for 7 days to limit disk usage.

### Promtail Configuration (promtail/config.yml)
```
server:
  http_listen_port: 9080

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: docker
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        filters:
          - name: label
            values: ["logging=promtail"]
        refresh_interval: 5s
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        regex: '/(.*)'
        target_label: 'container'
      - target_label: 'job'
        replacement: 'docker'
      - source_labels: ['__meta_docker_container_label_app']
        target_label: 'app'
```
Why this configuration?
- `docker_sd_configs` automatically discovers running containers via Docker socket.
- `relabel_configs` extracts container name into a label `container` and sets the `job` label.
- `pipeline_stages` includes `docker{}` to parse Docker’s JSON log format.
- No explicit file paths are needed because Promtails mounts `/var/lib/docker/containers` to read log files.

## Application Logging
Both applications (Python Flask and Go) were modified to output logs in JSON format.
### Python Flask (using custom JSONFormatter)
```
import json
import logging
from datetime import datetime

class JSONFormatter(logging.Formatter):
    def format(self, record):
        log_obj = {
            "timestamp": datetime.utcnow().isoformat() + "Z",
            "level": record.levelname,
            "message": record.getMessage(),
            "logger": record.name,
        }
        if record.exc_info:
            log_obj["exception"] = self.formatException(record.exc_info)
        return json.dumps(log_obj)

handler = logging.StreamHandler()
handler.setFormatter(JSONFormatter())
logger = logging.getLogger()
logger.addHandler(handler)
logger.setLevel(logging.INFO)
```
Resulting log entry:
```
{"timestamp": "2026-03-09T19:06:41.307Z", "level": "INFO", "message": "Request: GET / from 172.18.0.1", "logger": "__main__"}
```

### Go Application (standard log package)
Go app uses standard log.Printf which produces plain text. To keep structured logging, we could implement a JSON logger, but for this lab the logs are still collected. In a production environment, a proper JSON logger would be preferred.

## Dashboard
A simple dashboard was created in Grafana to visualize logs and metrics. The main panels:

### Panel 1: Logs Table (Logs visualization)
Query: `{app=~"devops-.*"}`
Explanation: Shows recent logs from all apps

### Panel 2: Request Rate (Time series graph)
Query: `sum by (app) (rate({app=~"devops-.*"} [1m]))`
Explanation: Shows logs per second by app

### Panel 3: Error Logs (Logs visualization)
Query: `{app=~"devops-.*"} | json | level="ERROR"`
Explanation: Shows only ERROR level logs

### Panel 4:      

Log Level Distribution (Stat or Pie chart)
Query: `sum by (level) (count_over_time({app=~"devops-.*"} | json [5m]))`
Explanation: Count logs by level (INFO, ERROR, etc.)

## Production Config
To make the stack production-ready, the following enhancements were applied:
### Resource Limits
Each service has memory and CPU limits to prevent resource exhaustion
### Security (Grafana)
- Anonymous access disabled (GF_AUTH_ANONYMOUS_ENABLED=false)
- Admin password stored in .env file (excluded from version control)
- Grafana requires login to view dashboards
### Health Checks
Health checks ensure services are responsive
### Log Retention
Loki configured to keep logs for 7 days (retention_period: 168h).

## Testing
### Health Check Status
```
docker compose ps
```
Expected output shows all services with healthy status (for those with healthchecks).

### Promtail Targets
Visit http://localhost:9080/targets. All application containers should be listed with status TRUE.

### Grafana Login
Access http://localhost:3000. You should see the login page, not the dashboard directly. Log in with credentials from .env.

### Query Logs in Explore
- Run a simple query: {app="devops-python"}
- Verify logs appear with JSON structure.

### Command-line Verification
```
curl -s http://localhost:3100/loki/api/v1/labels | jq .
```
Should return labels like app, container, job.
