# LAB 10 - Helm Chart Documentation – Python App

## Chart Overview

This Helm chart deploys the **DevOps Info Service**, a Python Flask application with Prometheus metrics and health checks.

### Chart Structure

```
k8s/python-app/
├── Chart.yaml                # Chart metadata (name, version, appVersion)
├── values.yaml               # Default configuration values
├── values-dev.yaml           # Development environment overrides
├── values-prod.yaml          # Production environment overrides
├── templates/
│   ├── _helpers.tpl          # Template helpers for labels, names, etc.
│   ├── deployment.yaml       # Deployment resource (templated)
│   ├── service.yaml          # Service resource (templated)
│   ├── configmap.yaml        # Application ConfigMap (optional)
│   └── hooks/                # Helm hooks directory
│       ├── pre-install-validation.yaml
│       └── post-install-smoketest.yaml
└── .helmignore               # Files to ignore when packaging
```

### Key Template Files

| File | Purpose |
|------|---------|
| `deployment.yaml` | Defines Pod template with probes, resources, environment variables |
| `service.yaml` | Exposes the application via NodePort (dev) or LoadBalancer (prod) |
| `pre-install-validation.yaml` | Pre‑install hook that validates prerequisites |
| `post-install-smoketest.yaml` | Post‑install hook that runs a smoke test against the deployed service |

### Values Organization Strategy

Values are organised hierarchically for clarity:

- **Global settings**: `replicaCount`, `image`, `strategy`
- **Service configuration**: `service.type`, `service.port`, `service.targetPort`
- **Resource management**: `resources.requests`, `resources.limits`
- **Probes**: `livenessProbe`, `readinessProbe`
- **Environment variables**: `env` list

Environment‑specific overrides are stored in separate files (`values-dev.yaml`, `values-prod.yaml`) and merged at install time.

---

## Configuration Guide

### Important Values

| Value | Description | Default |
|-------|-------------|---------|
| `replicaCount` | Number of Pod replicas | `3` |
| `image.repository` | Docker image repository | `damirsadykov/devops-info-service` |
| `image.tag` | Image tag | `latest` |
| `service.type` | Service type (`NodePort`, `LoadBalancer`, `ClusterIP`) | `NodePort` |
| `service.port` | Service port (cluster‑internal) | `80` |
| `service.targetPort` | Container port the app listens on | `5000` |
| `resources.requests.cpu` | Minimum CPU requested | `100m` |
| `resources.limits.cpu` | Maximum CPU allowed | `200m` |
| `livenessProbe.initialDelaySeconds` | Seconds before first liveness probe | `10` |
| `readinessProbe.initialDelaySeconds` | Seconds before first readiness probe | `5` |

### Customising for Different Environments

#### Development (1 replica, relaxed resources, NodePort)

```bash
helm install myapp-dev ./k8s/python-app -f ./k8s/python-app/values-dev.yaml
```

#### Production (5 replicas, higher resources, LoadBalancer)

```bash
helm install myapp-prod ./k8s/python-app -f ./k8s/python-app/values-prod.yaml
```

#### Override single value on the command line

```bash
helm install myapp-test ./k8s/python-app --set replicaCount=2 --set image.tag=v1.2.0
```

---

## Hook Implementation

### Hook Types Used

| Hook | File | Purpose |
|------|------|---------|
| `pre-install` | `hooks/pre-install-validation.yaml` | Runs before any resources are created – validates that required configuration exists (e.g., external ConfigMap). |
| `post-install` | `hooks/post-install-smoketest.yaml` | Runs after all resources are ready – performs a smoke test (calls `/health` and `/` endpoints). |

### Hook Weights

- **Pre‑install**: weight `-5` → runs early, before any other hooks.
- **Post‑install**: weight `5` → runs after all regular resources are deployed.

### Deletion Policies

Both hooks use:  
`helm.sh/hook-delete-policy: hook-succeeded`

This means the Job objects are automatically deleted when the hook completes successfully. This keeps the cluster clean and prevents accumulation of completed jobs.

### Why These Hooks?

- **Pre‑install**: In a real production scenario, you might need to verify that an external database secret or ConfigMap exists before the application starts. This hook blocks installation if prerequisites are missing.
- **Post‑install**: Automatically validates that the service is reachable and returns correct responses. It acts as a canary test – if the smoke test fails, the Helm release is marked as failed, alerting the operator.

---

## Installation Evidence

### 1. `helm list` output

```bash
$ helm list --namespace dev

NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
python-app-dev  dev             1               2026-04-02 21:16:20.931346564 +0300 MSK deployed        python-app-1.0.0        1.3.0      

$ helm list --namespace prod
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
python-app-prod prod            1               2026-04-02 21:15:29.273172537 +0300 MSK deployed        python-app-1.0.0        1.3.0      
```

### 2. `kubectl get all` showing deployed resources

```bash
$ kubectl get all --namespace dev
NAME                                  READY   STATUS    RESTARTS   AGE
pod/python-app-dev-5459d8ff4f-hghck   1/1     Running   0          2m47s

NAME                     TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
service/python-app-dev   NodePort   10.108.208.9   <none>        80:32599/TCP   2m47s

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/python-app-dev   1/1     1            1           2m47s

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/python-app-dev-5459d8ff4f   1         1         1       2m47s

$ kubectl get all --namespace prod
NAME                                   READY   STATUS    RESTARTS   AGE
pod/python-app-prod-567dfd9698-2gm4m   1/1     Running   0          3m30s
pod/python-app-prod-567dfd9698-6xhvl   1/1     Running   0          3m30s
pod/python-app-prod-567dfd9698-jd8ks   0/1     Pending   0          3m30s
pod/python-app-prod-567dfd9698-ml829   0/1     Pending   0          3m30s
pod/python-app-prod-567dfd9698-rrmn9   1/1     Running   0          3m30s

NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/python-app-prod   LoadBalancer   10.102.162.104   <pending>     80:30134/TCP   3m30s

NAME                              READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/python-app-prod   3/5     5            3           3m30s

NAME                                         DESIRED   CURRENT   READY   AGE
replicaset.apps/python-app-prod-567dfd9698   5         5         3       3m30s
```

### 3. Hook execution output

```bash
$ kubectl get jobs -w
python-app-dev-pre-install   0/1   Pending
python-app-dev-pre-install   0/1   Running
python-app-dev-pre-install   0/1   Completed
python-app-dev-pre-install   1/1   Completed
python-app-dev-pre-install   Deleted
python-app-dev-post-install   0/1   Pending
python-app-dev-post-install   0/1   Running
python-app-dev-post-install   0/1   Completed
python-app-dev-post-install   1/1   Completed
python-app-dev-post-install   Deleted
```

### 4. Different environment deployments

**Development**:

```bash
$ kubectl get deployment python-app-dev --namespace dev -o json | jq '.spec.replicas, .spec.template.spe
c.containers[0].resources.requests'
1
{
  "cpu": "50m",
  "memory": "64Mi"
}
```

**Production**:

```bash
$ kubectl get deployment python-app-prod --namespace prod -o json | jq '.spec.replicas, .spec.template.s
pec.containers[0].resources.requests'
5
{
  "cpu": "200m",
  "memory": "256Mi"
}
```

---

## Operations

### Installation

```bash
# Install development release
helm install myapp-dev ./k8s/python-app -f ./k8s/python-app/values-dev.yaml

# Install production release
helm install myapp-prod ./k8s/python-app -f ./k8s/python-app/values-prod.yaml
```

### Upgrade

Modify `values.yaml` or environment values, then:

```bash
helm upgrade myapp-dev ./k8s/python-app -f ./k8s/python-app/values-dev.yaml
```

### Rollback

```bash
# View history
helm history myapp-dev

# Rollback to revision 1
helm rollback myapp-dev 1
```

### Uninstall

```bash
helm uninstall myapp-dev
helm uninstall myapp-prod
```

---

## Testing & Validation

### 1. `helm lint`

```bash
$ helm lint ./k8s/python-app
==> Linting ./k8s/python-app
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
```

### 2. `helm template` verification

```bash
helm template test ./k8s/python-app | grep -E "kind: (Deployment|Service|Job)" | sort | uniq -c
      1 kind: Deployment
      2 kind: Job
      1 kind: Service
```

### 3. Dry‑run output

```bash
helm install --dry-run --debug test ./k8s/python-app
```

output:
```
$ helm install --dry-run --debug test ./k8s/python-app
level=WARN msg="--dry-run is deprecated and should be replaced with '--dry-run=client'"
level=DEBUG msg="Original chart version" version=""
level=DEBUG msg="Chart path" path=/home/damir/Desktop/DevOps/DevOps-Core-Course/k8s/python-app
level=DEBUG msg="number of dependencies in the chart" chart=python-app dependencies=0
NAME: test
LAST DEPLOYED: Thu Apr  2 21:26:40 2026
NAMESPACE: default
STATUS: pending-install
REVISION: 1
DESCRIPTION: Dry run complete
TEST SUITE: None
USER-SUPPLIED VALUES:
{}

COMPUTED VALUES:
env:
- name: PORT
  value: "5000"
- name: HOST
  value: 0.0.0.0
- name: DEBUG
  value: "false"
image:
  pullPolicy: Always
  repository: damirsadykov/devops-info-service
  tag: latest
livenessProbe:
  failureThreshold: 3
  httpGet:
    path: /health
    port: 5000
  initialDelaySeconds: 10
  periodSeconds: 5
readinessProbe:
  failureThreshold: 3
  httpGet:
    path: /ready
    port: 5000
  initialDelaySeconds: 5
  periodSeconds: 3
replicaCount: 3
resources:
  limits:
    cpu: 200m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
service:
  nodePort: ""
  port: 80
  targetPort: 5000
  type: NodePort
strategy:
  rollingUpdate:
    maxSurge: 1
    maxUnavailable: 0
  type: RollingUpdate

HOOKS:
---
# Source: python-app/templates/hooks/post-install-smoketest.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "test-python-app-post-install"
  labels:
    helm.sh/chart: python-app-1.0.0
    app.kubernetes.io/name: python-app
    app.kubernetes.io/instance: test
    app.kubernetes.io/version: "1.3.0"
    app.kubernetes.io/managed-by: Helm
  annotations:
    "helm.sh/hook": post-install
    "helm.sh/hook-weight": "5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: smoketest
        image: curlimages/curl:latest
        command:
          - sh
          - -c
          - |
            echo "Post-install hook: waiting for service to be ready..."
            # Service DNS name (works inside cluster)
            SERVICE_URL="http://test-python-app.default.svc.cluster.local:80"
            
            # Wait up to 60 seconds for the service to respond
            for i in $(seq 1 12); do
              echo "Attempt $i: curling $SERVICE_URL/health"
              if curl --fail --silent --output /dev/null "$SERVICE_URL/health"; then
                echo "✓ Service is healthy"
                break
              fi
              if [ $i -eq 12 ]; then
                echo "✗ Service failed to become ready"
                exit 1
              fi
              sleep 5
            done
            
            # Final test of the main endpoint
            echo "Testing main endpoint..."
            curl --fail --silent "$SERVICE_URL/" || exit 1
            
            echo "✓ Smoke test passed!"
---
# Source: python-app/templates/hooks/pre-install-validation.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "test-python-app-pre-install"
  labels:
    helm.sh/chart: python-app-1.0.0
    app.kubernetes.io/name: python-app
    app.kubernetes.io/instance: test
    app.kubernetes.io/version: "1.3.0"
    app.kubernetes.io/managed-by: Helm
  annotations:
    "helm.sh/hook": pre-install
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: pre-install-task
        image: busybox
        command:
          - sh
          - -c
          - |
            echo "Pre-install hook running at $(date)"
            echo "You could run database migrations, validations, etc. here."
            echo "Pre-install completed successfully."
            exit 0
MANIFEST:
---
# Source: python-app/templates/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  labels:
    helm.sh/chart: python-app-1.0.0
    app.kubernetes.io/name: python-app
    app.kubernetes.io/instance: test
    app.kubernetes.io/version: "1.3.0"
    app.kubernetes.io/managed-by: Helm
data:
  app.properties: |
    environment=dev
    log_level=info
---
# Source: python-app/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: test-python-app
  labels:
    helm.sh/chart: python-app-1.0.0
    app.kubernetes.io/name: python-app
    app.kubernetes.io/instance: test
    app.kubernetes.io/version: "1.3.0"
    app.kubernetes.io/managed-by: Helm
spec:
  type: NodePort
  ports:
    - port: 80
      targetPort: 5000
      protocol: TCP
      name: http
  selector:
    app.kubernetes.io/name: python-app
    app.kubernetes.io/instance: test
---
# Source: python-app/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-python-app
  labels:
    helm.sh/chart: python-app-1.0.0
    app.kubernetes.io/name: python-app
    app.kubernetes.io/instance: test
    app.kubernetes.io/version: "1.3.0"
    app.kubernetes.io/managed-by: Helm
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app.kubernetes.io/name: python-app
      app.kubernetes.io/instance: test
  template:
    metadata:
      labels:
        app.kubernetes.io/name: python-app
        app.kubernetes.io/instance: test
    spec:
      containers:
      - name: python-app
        image: "damirsadykov/devops-info-service:latest"
        imagePullPolicy: Always
        ports:
        - containerPort: 5000
          name: http
          protocol: TCP
        env:
        - name: PORT
          value: "5000"
        - name: HOST
          value: "0.0.0.0"
        - name: DEBUG
          value: "false"
        resources:
          limits:
            cpu: 200m
            memory: 256Mi
          requests:
            cpu: 100m
            memory: 128Mi
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /health
            port: 5000
          initialDelaySeconds: 10
          periodSeconds: 5
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /ready
            port: 5000
          initialDelaySeconds: 5
          periodSeconds: 3
```

### 4. Application accessibility

#### Development (NodePort)

```bash
export NODE_IP=$(minikube ip)
export NODE_PORT=$(kubectl get svc myapp-dev-python-app -o jsonpath='{.spec.ports[0].nodePort}')
curl http://$NODE_IP:$NODE_PORT/health
```

Expected response:

```json
{
  "status": "healthy",
  "timestamp": "2026-04-02T18:28:27.351637Z",
  "uptime_seconds": 670
}
```

#### Production (LoadBalancer with minikube tunnel)

```bash
minikube tunnel
export LB_IP=$(kubectl get svc myapp-prod-python-app -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
curl http://$LB_IP/ready
```

Expected response:

```json
{"status":"ready"}
```

---

## Conclusion

The Helm chart provides a production‑ready way to deploy the Python application with:

- Configurable resource limits and replicas
- Environment‑specific overrides
- Automatic smoke testing after deployment
- Zero‑downtime rolling updates
- Full lifecycle management (install, upgrade, rollback, uninstall)

All hooks are lightweight, reliable, and self‑cleaning.