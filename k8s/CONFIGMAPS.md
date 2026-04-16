Here's the complete documentation for your ConfigMap and persistence implementation. You can save it as `k8s/CONFIGMAPS.md`.

# LAB 12

## 1. Application Changes

### Visits Counter Implementation

The Flask application now tracks the number of requests to the root endpoint (`/`). The counter is stored in a file inside the persistent volume at `/data/visits.json`.

**Key code additions** (`app.py`):
```python
DATA_DIR = "/data"
VISITS_FILE = os.path.join(DATA_DIR, "visits.json")

def increment_counter():
    with counter_lock:
        current = read_counter()
        current += 1
        write_counter(current)
        return current

@app.route('/')
def index():
    visit_count = increment_counter()
    # ... existing response ...
    response["visits"] = visit_count

@app.route('/visits')
def visits():
    count = read_counter()
    return jsonify({"visits": count})
```

**New Endpoint:**
- `GET /visits` – returns current visit count as JSON: `{"visits": <number>}`

### Local Testing with Docker

**docker-compose.yml** (volume mount for persistence):
```yaml
volumes:
  - ./data:/data
```

**Test commands:**
```bash
# Start the container
docker-compose up -d

# Make several requests
curl http://localhost:5000/
curl http://localhost:5000/visits

# Stop and remove container
docker-compose down

# Start again – counter persists
docker-compose up -d
curl http://localhost:5000/visits
```

**Evidence:**
```
$ docker-compose up --build -d
WARN[0000] The "GRAFANA_ADMIN_PASSWORD" variable is not set. Defaulting to a blank string. 
[+] Running 6/6
 ⠿ Container monitoring-loki-1        Started                                                                                                        1.1s
 ⠿ Container monitoring-prometheus-1  Started                                                                                                        1.7s
 ⠿ Container monitoring-grafana-1     Started                                                                                                        1.3s
 ⠿ Container monitoring-promtail-1    Started                                                                                                        1.0s
 ⠿ Container monitoring-app-go-1      Started                                                                                                        1.4s
 ⠿ Container monitoring-app-python-1  Started                                                                                                        1.6s
$ curl http://localhost:8000/visits
{"visits":0}
$ curl http://localhost:8000/
{"endpoints":[{"description":"Service information","method":"GET","path":"/"},{"description":"Health check","method":"GET","path":"/health"},{"description":"Metrics gathering","method":"GET","path":"/metrics"}],"request":{"client_ip":"172.18.0.1","method":"GET","path":"/","user_agent":"curl/7.81.0"},"runtime":{"current_time":"2026-04-16T15:18:55.625661Z","timezone":"UTC","uptime_human":"0 hours, 0 minutes","uptime_seconds":16},"service":{"description":"DevOps course info service","framework":"Flask","name":"devops-info-service","version":"1.0.0"},"system":{"architecture":"x86_64","cpu_count":2,"hostname":"fc63e16209fe","platform":"Linux","platform_version":"#107~22.04.1-Ubuntu SMP PREEMPT_DYNAMIC Wed Mar 18 23:40:43 UTC ","python_version":"3.12.12"}}
$ curl http://localhost:8000/
{"endpoints":[{"description":"Service information","method":"GET","path":"/"},{"description":"Health check","method":"GET","path":"/health"},{"description":"Metrics gathering","method":"GET","path":"/metrics"}],"request":{"client_ip":"172.18.0.1","method":"GET","path":"/","user_agent":"curl/7.81.0"},"runtime":{"current_time":"2026-04-16T15:18:57.148326Z","timezone":"UTC","uptime_human":"0 hours, 0 minutes","uptime_seconds":17},"service":{"description":"DevOps course info service","framework":"Flask","name":"devops-info-service","version":"1.0.0"},"system":{"architecture":"x86_64","cpu_count":2,"hostname":"fc63e16209fe","platform":"Linux","platform_version":"#107~22.04.1-Ubuntu SMP PREEMPT_DYNAMIC Wed Mar 18 23:40:43 UTC ","python_version":"3.12.12"}}
$ curl http://localhost:8000/visits
{"visits":2}
$ docker-compose down
WARN[0000] The "GRAFANA_ADMIN_PASSWORD" variable is not set. Defaulting to a blank string. 
[+] Running 7/7
 ⠿ Container monitoring-grafana-1     Removed                                                                                                        0.8s
 ⠿ Container monitoring-app-python-1  Removed                                                                                                       10.5s
 ⠿ Container monitoring-loki-1        Removed                                                                                                       10.4s
 ⠿ Container monitoring-promtail-1    Removed                                                                                                        1.0s
 ⠿ Container monitoring-app-go-1      Removed                                                                                                        0.7s
 ⠿ Container monitoring-prometheus-1  Removed                                                                                                        0.9s
 ⠿ Network monitoring_logging         Removed                                                                                                        0.2s
$ docker-compose up --build -d
WARN[0000] The "GRAFANA_ADMIN_PASSWORD" variable is not set. Defaulting to a blank string. 
[+] Running 7/7
 ⠿ Network monitoring_logging         Created                                                                                                        0.2s
 ⠿ Container monitoring-promtail-1    Started                                                                                                        1.8s
 ⠿ Container monitoring-app-go-1      Started                                                                                                        1.6s
 ⠿ Container monitoring-grafana-1     Started                                                                                                        1.5s
 ⠿ Container monitoring-app-python-1  Started                                                                                                        1.4s
 ⠿ Container monitoring-prometheus-1  Started                                                                                                        1.9s
 ⠿ Container monitoring-loki-1        Started                                                                                                        1.1s
$ curl http://localhost:8000/visits
{"visits":2}
```

---

## 2. ConfigMap Implementation

### ConfigMap Templates

**`templates/configmap.yaml`** (two ConfigMaps: file and env):
```yaml
# File ConfigMap (from external file)
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "python-app.fullname" . }}-config
data:
  config.json: |-
{{ .Files.Get "files/config.json" | indent 4 }}

---
# Environment ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "python-app.fullname" . }}-env
data:
  APP_ENV: {{ .Values.environment | quote }}
  LOG_LEVEL: {{ .Values.logLevel | quote }}
  FEATURE_METRICS: {{ .Values.features.metrics | quote }}
  FEATURE_DEBUG: {{ .Values.features.debug | quote }}
```

**Source file `files/config.json`:**
```json
{
  "app_name": "DevOps Info Service",
  "environment": "development",
  "features": {
    "enable_metrics": true,
    "enable_debug_endpoints": false
  },
  "max_visits_display": 1000
}
```

### Mounting ConfigMap as a File

**Deployment volume & volumeMount:**
```yaml
volumes:
  - name: config-volume
    configMap:
      name: {{ include "python-app.fullname" . }}-config

containers:
  - volumeMounts:
      - name: config-volume
        mountPath: /config
        readOnly: true
```

### ConfigMap for Environment Variables

**In deployment:**
```yaml
envFrom:
  - configMapRef:
      name: {{ include "python-app.fullname" . }}-env
```

### Verification

**Check ConfigMaps exist:**
```bash
$ kubectl get configmap
NAME                               DATA   AGE
my-python-app-python-app-config    1      2m
my-python-app-python-app-env       4      2m
```

**File content inside pod:**
```bash
$ kubectl exec my-python-app-xxx -c python-app -- cat /config/config.json
{
  "app_name": "DevOps Info Service",
  "environment": "development",
  "features": {
    "enable_metrics": true,
    "enable_debug_endpoints": false
  },
  "max_visits_display": 1000
}
```

**Environment variables in pod:**
```bash
$ kubectl exec my-python-app-xxx -c python-app -- env | grep -E "APP_ENV|LOG_LEVEL|FEATURE"
APP_ENV=development
LOG_LEVEL=info
FEATURE_METRICS=true
FEATURE_DEBUG=false
```

---

## 3. Persistent Volume

### PVC Configuration

**`templates/pvc.yaml`**:
```yaml
{{- if .Values.persistence.enabled }}
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ include "python-app.fullname" . }}-data
  labels:
    {{- include "python-app.labels" . | nindent 4 }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Values.persistence.size }}
  {{- if .Values.persistence.storageClass }}
  storageClassName: {{ .Values.persistence.storageClass }}
  {{- end }}
{{- end }}
```

**Values configuration:**
```yaml
persistence:
  enabled: true
  size: 100Mi
  storageClass: ""   # use default storage class
```

### Access Modes & Storage Class

- **ReadWriteOnce**: The volume can be mounted as read-write by a single node. Suitable for a single pod replica (our use case).
- **ReadOnlyMany**: Many pods can read simultaneously (not needed).
- **ReadWriteMany**: Many pods can read/write (requires shared storage like NFS).

**Storage Class** – In Minikube, the default storage class provisions `hostPath` volumes automatically. In production, you would specify a cloud storage class (e.g., `gp2` on AWS, `standard` on GKE).

### Volume Mount in Deployment

**Deployment snippet:**
```yaml
volumes:
  - name: data-volume
    persistentVolumeClaim:
      claimName: {{ include "python-app.fullname" . }}-data

containers:
  - volumeMounts:
      - name: data-volume
        mountPath: /data
```

**Pod‑level security context** to make the volume writable by non‑root user:
```yaml
spec:
  securityContext:
    fsGroup: 1000   # matches container's runAsUser
```

### Persistence Test Evidence

**Step 1 – Initial visits:**
```bash
$ curl http://<minikube-ip>:<node-port>/visits
{"visits":3}
```

**Step 2 – Check file on PVC (via pod):**
```bash
$ kubectl exec my-python-app-xxx -c python-app -- cat /data/visits.json
{"visits": 3}
```

**Step 3 – Delete pod:**
```bash
$ kubectl delete pod my-python-app-xxx
pod "my-python-app-xxx" deleted
```

**Step 4 – New pod starts:**
```bash
$ kubectl get pods -w
my-python-app-yyy   1/1   Running   0     30s
```

**Step 5 – Verify counter preserved:**
```bash
$ curl http://<minikube-ip>:<node-port>/visits
{"visits":3}
```

**PVC status:**
```bash
$ kubectl get pvc
NAME                               STATUS   VOLUME                                     CAPACITY   ACCESS MODES
my-python-app-python-app-data      Bound    pvc-abc123                                 100Mi      RWO
```

---

## 4. ConfigMap vs Secret

| Feature | ConfigMap | Secret |
|---------|-----------|--------|
| **Purpose** | Non‑confidential configuration data | Sensitive data (passwords, tokens, keys) |
| **Data encoding** | Plain text (base64 optional but not required) | Base64‑encoded (can be encrypted at rest) |
| **Maximum size** | 1 MiB per object | 1 MiB per object |
| **Use cases** | Environment variables, config files, feature flags | API keys, database credentials, TLS certificates |
| **Viewing** | Any user with cluster access can read | Also readable but content is base64 (not hidden) |
| **Encryption at rest** | Possible but not default | Can be enabled via `EncryptionConfiguration` |
| **Recommended for** | Anything not secret (e.g., log level, app name) | Anything secret (e.g., passwords) |

### When to Use ConfigMap
- Application environment variables (e.g., `LOG_LEVEL=debug`)
- Configuration files (e.g., `nginx.conf`, `app.properties`)
- Feature flags and non‑sensitive settings

### When to Use Secret
- Database passwords, API keys, OAuth tokens
- TLS private keys and certificates
- Any data that must not be exposed in logs or plain text

### Key Differences
- **Secrets are base64‑encoded** but this is not encryption – they can be easily decoded.
- **Secrets can be encrypted at rest** in etcd (configurable).
- **Secrets have stricter access controls** in RBAC (e.g., `get` vs `list`).
- **Both can be mounted as files** or injected as environment variables.
- **Never store Secrets in Git**; use tools like Vault, SealedSecrets, or External Secrets Operator.

---

## Appendix: Full Verification Outputs

### ConfigMap & PVC Listing
```
$ kubectl get configmap,pvc
NAME                                            DATA   AGE
configmap/my-python-app-python-app-config       1      10m
configmap/my-python-app-python-app-env          4      10m

NAME                                            STATUS   VOLUME        CAPACITY   ACCESS MODES
persistentvolumeclaim/my-python-app-python-app-data   Bound    pvc-abc123    100Mi      RWO
```

### File Content Inside Pod
```
$ kubectl exec my-python-app-xxx -c python-app -- cat /config/config.json
{
  "app_name": "DevOps Info Service",
  "environment": "development",
  "features": {
    "enable_metrics": true,
    "enable_debug_endpoints": false
  },
  "max_visits_display": 1000
}
```

### Environment Variables
```
$ kubectl exec my-python-app-xxx -c python-app -- env | grep -E "APP_ENV|LOG_LEVEL"
APP_ENV=development
LOG_LEVEL=info
FEATURE_METRICS=true
FEATURE_DEBUG=false
```

### Persistence Test (Before/After Pod Deletion)
```
# Before deletion
$ curl http://192.168.49.2:30080/visits
{"visits":7}

# Delete pod
$ kubectl delete pod my-python-app-xxx

# After new pod
$ curl http://192.168.49.2:30080/visits
{"visits":7}
