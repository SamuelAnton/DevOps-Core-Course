# Lab 8

## 1. Architecture

The monitoring stack follows the **pull‑based** Prometheus model. Below is a high‑level flow of metrics:

```
┌─────────────┐       scrape       ┌─────────────┐      query      ┌─────────────┐
│   Python    │ ◄───────────────── │ Prometheus  │ ◄────────────── │   Grafana   │
│     App     │   /metrics (15s)   │   Server    │                 │  Dashboard  │
└─────────────┘                    └─────────────┘                 └─────────────┘
       │                                    │
       │                                    │
       └────── logs ──────► Loki ──────────┘
                        (via Promtail)
```

- **Python App** exposes custom and HTTP metrics at `/metrics`.
- **Prometheus** scrapes the app every 15 seconds, stores time‑series data locally, and applies retention policies.
- **Grafana** queries Prometheus to visualise metrics in dashboards.
- **Loki** (from Lab 7) collects logs, which complement metrics for full observability.

All components run in Docker Compose, share the `logging` network, and have health checks and resource limits.

---

## 2. Application Instrumentation

I instrumented the Python Flask application using the `prometheus_client` library. The following metrics were added:

| Metric                           | Type      | Labels                          | Purpose |
|----------------------------------|-----------|---------------------------------|---------|
| `http_requests_total`            | Counter   | `method`, `endpoint`, `status` | Count total HTTP requests, split by method, endpoint, and response status. |
| `http_request_duration_seconds`  | Histogram | `method`, `endpoint`            | Measure request duration distribution to detect slow endpoints. |
| `http_requests_in_progress`      | Gauge     | (none)                          | Track number of concurrent requests – useful for capacity planning. |
| `system_info_collection_seconds` | Histogram | (none)                          | Measure time spent collecting system information in the `/` endpoint. |

**Why these metrics?**
- **Counter** for request rate and error rate (RED method: Rate, Errors, Duration).
- **Histogram** for latency percentiles (e.g., p95).
- **Gauge** for concurrency, which helps spot bottlenecks.
- **Business metric** (`system_info_collection_seconds`) monitors a specific operation that could become a performance issue.

Labels are kept low‑cardinality (`method`, `endpoint`, `status`) to avoid bloating the time‑series database.

---

## 3. Prometheus Configuration

**File:** `monitoring/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

storage:
  tsdb:
    retention.time: 15d      # Keep data for 15 days
    retention.size: 10GB     # Limit TSDB size

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'app'
    static_configs:
      - targets: ['app-python:5000']
    metrics_path: /metrics

  - job_name: 'loki'
    static_configs:
      - targets: ['loki:3100']
    metrics_path: /metrics

  - job_name: 'grafana'
    static_configs:
      - targets: ['grafana:3000']
    metrics_path: /metrics
```

**Explanation:**
- **Scrape interval:** 15 seconds – a good balance between freshness and resource usage.
- **Retention:** 15 days or 10GB (whichever is reached first) to manage disk space.
- **Targets:** Prometheus itself, the Python app, Loki, and Grafana. All are reachable via Docker service names on the `logging` network.

---

## 4. Dashboard Walkthrough

We created a custom dashboard in Grafana with seven panels, each visualising a different aspect of the application.

### Panel 1: Request Rate (per endpoint)
- **Query:** `sum(rate(http_requests_total[5m])) by (endpoint)`
- **Visualisation:** Time series
- **Legend:** `{{endpoint}}`
- **Unit:** requests/sec (reqps)
- **Purpose:** Shows traffic distribution across endpoints.

### Panel 2: Error Rate (5xx)
- **Query:** `sum(rate(http_requests_total{status=~"5.."}[5m]))`
- **Visualisation:** Time series
- **Unit:** reqps
- **Purpose:** Alerts when server errors increase.

### Panel 3: Request Duration p95
- **Query:** `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, endpoint))`
- **Visualisation:** Time series
- **Legend:** `{{endpoint}} p95`
- **Unit:** seconds
- **Purpose:** Monitors latency experienced by 95% of users.

### Panel 4: Request Duration Heatmap
- **Query:** `rate(http_request_duration_seconds_bucket[5m])`
- **Visualisation:** Heatmap
- **Y‑Axis Unit:** seconds
- **Purpose:** Visualises latency distribution over time; reveals whether latency spikes affect all requests or only a few.

### Panel 5: Active Requests
- **Query:** `sum(http_requests_in_progress)`
- **Visualisation:** Stat or Time series
- **Unit:** short
- **Purpose:** Shows current concurrency; helps detect traffic surges or stuck requests.

### Panel 6: Status Code Distribution
- **Query:** `sum by (status) (rate(http_requests_total[5m]))`
- **Visualisation:** Pie chart or Bar gauge
- **Legend:** `{{status}}`
- **Unit:** reqps
- **Purpose:** Quickly see the proportion of 2xx, 4xx, 5xx responses.

### Panel 7: Uptime
- **Query:** `up{job="app"}`
- **Visualisation:** Stat
- **Unit:** none (1 = up, 0 = down)
- **Purpose:** Simple service health indicator.

All panels use `$__interval` and `$__rate_interval` to adapt to the selected time range.

---

## 5. PromQL Examples

1. **Overall request rate**  
   `rate(http_requests_total[5m])`  
   *Returns requests per second for every time series.*

2. **Error rate percentage**  
   `(sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))) * 100`  
   *Shows the percentage of 5xx errors.*

3. **95th percentile latency for the root endpoint**  
   `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{endpoint="/"}[5m])) by (le))`  
   *Focuses on a specific endpoint.*

4. **Top 3 slowest endpoints**  
   `topk(3, avg by (endpoint) (rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])))`  
   *Computes average latency per endpoint and picks the highest three.*

5. **Concurrent requests over the last hour**  
   `avg_over_time(http_requests_in_progress[1h])`  
   *Smooths the gauge to show typical concurrency.*

---

## 6. Production Setup

### Health Checks
Every service in `docker-compose.yml` includes a health check to ensure it's truly ready:

- **Prometheus:** `wget --spider http://localhost:9090/-/healthy`
- **Python App:** `curl -f http://localhost:5000/health` (or `wget --spider` if curl is missing)
- **Grafana:** `wget --spider http://localhost:3000/api/health`
- **Loki:** `wget --spider http://localhost:3100/ready`

Checks run every 10 seconds; after 5 retries the container is marked unhealthy.

### Resource Limits
All services have CPU and memory limits (using `deploy.resources` for compatibility with Docker Compose v3):

| Service   | CPU Limit | Memory Limit |
|-----------|-----------|--------------|
| Prometheus| 1.0       | 1GB          |
| Loki      | 1.0       | 1GB          |
| Grafana   | 0.5       | 512MB        |
| Python App| 0.5       | 256MB        |

### Data Retention
- **Prometheus:** 15 days or 10GB (command‑line flags `--storage.tsdb.retention.time=15d` and `--storage.tsdb.retention.size=10GB`).
- **Loki:** 7 days (configured in `loki/config.yml` with `retention_period: 168h`).
- **Grafana:** Dashboards are stored in a persistent volume and kept indefinitely.

### Persistent Volumes
Three named volumes guarantee data survives restarts:
- `prometheus-data` (mounted at `/prometheus`)
- `loki-data` (mounted at `/loki`)
- `grafana-data` (mounted at `/var/lib/grafana`)

**Test:** After creating a dashboard, running `docker compose down` and `up -d`, the dashboard reappears – confirming persistence.

---

## 7. Testing Results
All in "screenshots/lab 8" folder

---

## 8. Metrics vs. Logs (Lab 7)

| Aspect         | Metrics (Prometheus)                          | Logs (Loki)                                   |
|----------------|-----------------------------------------------|-----------------------------------------------|
| **Data type**  | Structured numbers (counters, gauges, histograms) | Unstructured text (with optional labels)      |
| **Purpose**    | Trend analysis, alerting, SLI/SLO monitoring  | Debugging, detailed event inspection, audit   |
| **Query language** | PromQL (aggregation, rate, quantiles)     | LogQL (filtering, pattern matching)           |
| **Storage**    | Time‑series database, efficient for numeric data | Object store, compressed chunks                |
| **Example use**| “Request error rate exceeded 5%”              | “Show me all 500 errors for user X”           |
| **Retention**  | Short‑term (days) – high cardinality          | Longer‑term (weeks/months) – raw data         |

**When to use each:**
- **Metrics** for dashboards, alerts, and capacity planning.
- **Logs** for root‑cause analysis and compliance.

Together they provide complete observability.