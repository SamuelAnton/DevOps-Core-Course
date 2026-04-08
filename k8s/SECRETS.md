# LAB 11

## Kubernetes Secrets
### Secret creation and demonstration
```
$ kubectl create secret generic app-credentials \
  --from-literal=username=password \
  --from-literal=password=usernmae
secret/app-credentials created

$ kubectl get secret app-credentials -o yaml
apiVersion: v1
data:
  password: dXNlcm5tYWU=
  username: cGFzc3dvcmQ=
kind: Secret
metadata:
  creationTimestamp: "2026-04-07T15:46:39Z"
  name: app-credentials
  namespace: default
  resourceVersion: "12198"
  uid: 4fa41251-fa63-4ed6-9f1c-b840d2f96cb8
type: Opaque

$ echo "dXNlcm5tYWU=" | base64 -d
usernmaedamir
$ echo "cGFzc3dvcmQ=" | base64 -d
passworddamir 
```

### Base64 encoding vs encryption
- **Base64 encoding** is reversible without a key. It only converts binary data to ASCII for safe transmission in YAML/JSON.
- **Encryption** (e.g., AES) requires a secret key and is not reversible without that key.
- Kubernetes Secrets are base64‑encoded by default. To enable encryption at rest, you must configure an **EncryptionConfiguration** or use an external secrets manager like Vault.


## Helm Secret Integration

### Chart Structure

The Helm chart includes a `templates/secrets.yaml` file that defines the Secret resource. This file uses Helm templating to reference values from `values.yaml`.
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: {{ include "python-app.fullname" . }}-secret
  labels:
    {{- include "python-app.labels" . | nindent 4 }}
type: Opaque
stringData:
  API_KEY: {{ .Values.secrets.apiKey | default "change-me" }}
  DB_PASSWORD: {{ .Values.secrets.dbPassword | default "change-me" }}
  ADMIN_PASSWORD: {{ .Values.secrets.adminPassword | default "change-me" }}
```

The `values.yaml` contains placeholder defaults (never commit real secrets):

```yaml
secrets:
  apiKey: "dev-key"
  dbPassword: "dev-db-pass"
  adminPassword: "dev-admin-pass"
```


### Consuming Secrets in Deployment

The `deployment.yaml` injects the secret keys as environment variables using `secretKeyRef`:

```yaml
env:
  - name: API_KEY
    valueFrom:
      secretKeyRef:
        name: {{ include "python-app.secretName" . }}
        key: API_KEY
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: {{ include "python-app.secretName" . }}
        key: DB_PASSWORD
  - name: ADMIN_PASSWORD
    valueFrom:
      secretKeyRef:
        name: {{ include "python-app.secretName" . }}
        key: ADMIN_PASSWORD
```

### Verification (Environment Variables Inside Pod)

```bash
$ kubectl exec my-python-app-64d5cbf966-zqqnm -c python-app -- env | grep -E "API_KEY|DB_PASSWORD|ADMIN_PASSWORD"
API_KEY=dev-key
DB_PASSWORD=dev-db-pass
ADMIN_PASSWORD=dev-admin-pass
```

## Resource Management

### Resource Limits Configuration

In `values.yaml`, we define CPU and memory requests and limits:

```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "200m"
```

The deployment applies these values using:

```yaml
resources:
  {{- toYaml .Values.resources | nindent 10 }}
```

### Explanation: Requests vs Limits

| Concept | Purpose | Example |
|---------|---------|---------|
| **Requests** | Minimum resources guaranteed to the container. The scheduler uses this to place pods. | `cpu: 100m` (0.1 core), `memory: 128Mi` |
| **Limits** | Maximum resources the container can use. If exceeded, the container may be throttled (CPU) or terminated (memory). | `cpu: 200m`, `memory: 256Mi` |

**Why both are important:**  
- Setting **requests** ensures the pod gets enough resources to run.  
- Setting **limits** prevents a single pod from starving other pods on the node.  
- A well‑tuned configuration balances performance and reliability.

### Choosing Appropriate Values

- **Start with monitoring:** Observe resource usage during normal load (using `kubectl top pod`, Prometheus, etc.).
- **Requests:** Set to average baseline usage + small buffer.
- **Limits:** Set to peak usage + safety margin (e.g., 1.5× requests).
- For our Python Flask app with Prometheus metrics, we chose 100m CPU / 128Mi memory as requests, and 200m CPU / 256Mi memory as limits.

## Vault Integration

### Vault installation verification
```
$ helm install vault hashicorp/vault   --set "server.dev.enabled=true"   --set "injector.enabled=true"
NAME: vault
LAST DEPLOYED: Wed Apr  8 18:50:20 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
NOTES:
Thank you for installing HashiCorp Vault!

Now that you have deployed Vault, you should look over the docs on using
Vault with Kubernetes available here:

https://developer.hashicorp.com/vault/docs


Your release is named vault. To learn more about the release, try:

  $ helm status vault
  $ helm get manifest vault

$ kubectl get pods
NAME                                    READY   STATUS              RESTARTS   AGE
vault-0                                 0/1     ContainerCreating   0          33s
vault-agent-injector-848dd747d7-45xhm   1/1     Running             0          34s
```

### Policy and Role Configuration

**Policy (`myapp-policy`):** Grants read access to the secret path.

```hcl
path "secret/data/myapp/*" {
  capabilities = ["read", "list"]
}
path "secret/metadata/myapp/*" {
  capabilities = ["read", "list"]
}
```

**Role (`myapp-role`):** Binds the policy to the Kubernetes service account `my-python-app-sa`.

```bash
$ vault read auth/kubernetes/role/myapp-role
Key                                         Value
---                                         -----
bound_service_account_names                 [my-python-app-sa]
bound_service_account_namespaces            [default]
policies                                    [myapp-policy]
token_ttl                                   24h
```

### Proof of Secret Injection

The Vault Agent sidecar injects secrets as a file inside the pod at `/vault/secrets/config`.

**Check that the file exists:**

```bash
$ kubectl exec my-python-app-64d5cbf966-zqqnm -c python-app -- ls -la /vault/secrets/
total 4
drwxrwxrwt 3 root root      60 Apr  8 17:04 .
drwxr-xr-x 3 root root    4096 Apr  8 17:04 ..
drwxr-xr-x 3  100 appuser   60 Apr  8 17:04 vault
```

**Display the secret content:**

```bash
$ kubectl exec my-python-app-64d5cbf966-zqqnm -c python-app -- cat /vault/secrets/config
username=admin
password=secret123
```

### Explanation of the Sidecar Injection Pattern

- **Vault Agent Injector** is a mutating admission webhook that modifies pod definitions at creation time.
- When a pod has the annotation `vault.hashicorp.com/agent-inject: "true"`, the injector adds a **sidecar container** (the Vault Agent) to the pod.
- The sidecar authenticates to Vault using the pod’s service account token and retrieves the requested secrets.
- Secrets are written to a shared emptyDir volume mounted in both the sidecar and the application container.
- The application reads secrets from a file (or environment variables, if templated) without ever knowing the Vault address or token.

**Benefits:**  
- No secrets in environment variables (which can be leaked via `/proc` or debugging tools).  
- Secrets are never written to disk except in `tmpfs` (memory).  
- Automatic renewal of tokens and secrets.

## 5. Security Analysis

### Comparison: Kubernetes Secrets vs HashiCorp Vault

| Feature | Kubernetes Secrets | HashiCorp Vault |
|---------|--------------------|------------------|
| **Storage** | etcd (base64‑encoded, can be encrypted at rest) | Encrypted backend (integrated storage, Raft, etc.) |
| **Access Control** | RBAC (coarse) | Fine‑grained policies, path‑based ACLs |
| **Dynamic Secrets** | No | Yes (e.g., database credentials on‑demand) |
| **Audit Logging** | Limited (API server logs) | Detailed audit logs with request/response |
| **Secret Rotation** | Manual | Automatic (with leases and TTL) |
| **Integration** | Native to Kubernetes | Many backends (cloud, databases, PKI) |
| **Operational Overhead** | Minimal | Additional cluster to manage |

### When to Use Each Approach

**Use Kubernetes Secrets when:**
- You need a simple, built‑in solution for non‑critical secrets (e.g., API keys for internal services).
- Your compliance requirements are low.
- You are already using Helm and want to keep the stack minimal.

**Use Vault when:**
- You have **high‑value secrets** (e.g., database passwords, TLS certificates, cloud credentials).
- You need **dynamic secrets** (e.g., per‑pod database credentials that expire).
- You require **audit logging** of who accessed which secret.
- You need **automatic secret rotation** without redeploying applications.
- You operate in a regulated environment (PCI‑DSS, HIPAA, SOC2).

### Production Recommendations

1. **Never store secrets in Git** – use `--set` with CI/CD, external secrets operator, or Vault.
2. **Enable encryption at rest** for Kubernetes Secrets (if using them) via `EncryptionConfiguration`.
3. **Use Vault’s Kubernetes auth method** with short‑lived service account tokens.
4. **Restrict access** with least‑privilege policies (e.g., `read` only on specific paths).
5. **Audit all secret access** – enable Vault’s audit log and ship to a SIEM.
6. **For high‑scale environments**, prefer Vault Agent sidecar over direct API calls from the application.
7. **Consider using the External Secrets Operator** to sync Vault secrets into Kubernetes Secrets (hybrid approach).
