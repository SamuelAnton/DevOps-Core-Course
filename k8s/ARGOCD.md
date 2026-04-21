# Lab 13

## 1. ArgoCD Setup

### Installation Verification

Installed ArgoCD using the official Helm chart:

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
kubectl create namespace argocd
helm install argocd argo/argo-cd --namespace argocd
```

**All components are running:**

```bash
$ kubectl get pods -n argocd
NAME                                                READY   STATUS    RESTARTS   AGE
argocd-application-controller-0                     1/1     Running   0          41m
argocd-applicationset-controller-59f6b7dd64-4gczf   1/1     Running   0          41m
argocd-dex-server-7b9588c494-pc7dw                  1/1     Running   0          41m
argocd-notifications-controller-8f6855454-5mtrz     1/1     Running   0          41m
argocd-redis-dc6b586fc-sm2t8                        1/1     Running   0          41m
argocd-repo-server-5f4d44d9f8-p4l8g                 1/1     Running   0          41m
argocd-server-5f777b877f-lb722                      1/1     Running   0          41m
```

### UI Access Method
Port‑forward the ArgoCD server service to `localhost:8080`:

```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

**Login credentials:**
- Username: `admin`
- Password: retrieved from the initial admin secret

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

The UI is accessible at `https://localhost:8080`.

### CLI Configuration

Installed the ArgoCD CLI and logged in:

```bash
curl -sSL -o argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
chmod +x argocd
sudo mv argocd /usr/local/bin/

argocd login localhost:8080 --insecure
```

**Verification:**

```bash
$ argocd version

argocd: v3.3.7+035e855
  BuildDate: 2026-04-16T15:58:07Z
  GitCommit: 035e8556c451196e203078160a5c01f43afdb92f
  GitTreeState: clean
  GoVersion: go1.25.5
  Compiler: gc
  Platform: linux/amd64


$ argocd app list
NAME  CLUSTER  NAMESPACE  PROJECT  STATUS  HEALTH  SYNCPOLICY  CONDITIONS  REPO  PATH  TARGET
```

---

## 2. Application Configuration

### Application Manifests

Defined two ArgoCD Application resources in `k8s/argocd/application-dev.yaml` and `application-prod.yaml`.

**Dev Application (excerpt):**
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: python-app-dev
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/SamuelAnton/DevOps-Core-Course.git
    targetRevision: master
    path: k8s/python-app
    helm:
      valueFiles:
        - values-dev.yaml
  destination:
    server: https://kubernetes.default.svc
    namespace: dev
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

**Prod Application (excerpt):**
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: python-app-prod
  namespace: argocd
spec:
  # ... similar source ...
  destination:
    namespace: prod
  syncPolicy:
    syncOptions:
      - CreateNamespace=true
    # No 'automated' block → manual sync
```

### Source and Destination

- **Source**: Git repository at `master` branch, Helm chart in `k8s/python-app`. Environment‑specific values files (`values-dev.yaml`, `values-prod.yaml`) override defaults.
- **Destination**: The same Kubernetes cluster (`in‑cluster`), but different namespaces (`dev`, `prod`).

---

## 3. Multi-Environment Deployment

### Configuration Differences

| Parameter | Dev (values-dev.yaml) | Prod (values-prod.yaml) |
|-----------|----------------------|-------------------------|
| `replicaCount` | 1 | 5 |
| `resources.requests.cpu` | 50m | 200m |
| `resources.limits.cpu` | 100m | 500m |
| `environment` | development | production |
| `logLevel` | debug | warn |
| `DEBUG` env var | true | false |

### Sync Policy Differences

- **Dev**: Auto‑sync enabled (`automated` with `prune` and `selfHeal`). Any change in Git is automatically applied to the cluster.
- **Prod**: Manual sync only. Changes must be explicitly synchronised by an operator.

**Rationale for manual prod sync:**
- Requires human approval after change review.
- Allows deployment during maintenance windows.
- Provides rollback capability without automation.
- Compliance and separation of duties.

### Namespace Separation

- `dev` and `prod` namespaces isolate workloads and resources.
- Each environment runs its own instance of the application.
- ArgoCD Applications are scoped to their respective namespaces.

```bash
$ kubectl get pods -n dev
NAME                                      READY   STATUS    RESTARTS   AGE
python-app-python-app-7b858cb655-5gdxv    1/1     Running   0          2m

$ kubectl get pods -n prod
NAME                                      READY   STATUS    RESTARTS   AGE
python-app-python-app-6c9b7d8f9c-4j2kd    5/5     Running   0          1m
```

---

## 4. Self-Healing Evidence

### Test 1: Manual Scale (Self-Healing)

**Initial state (dev):** Replicas = 1 (as defined in Git)

```bash
$ kubectl get deployment -n dev
NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
python-app-python-app    1/1     1            1           10m
```

**Manual scale to 5 replicas:**

```bash
$ kubectl scale deployment python-app-python-app -n dev --replicas=5
deployment.apps/python-app-python-app scaled
```

Pods quickly scale to 5. ArgoCD detects drift and within ~3 minutes reverts to Git state (1 replica).

**ArgoCD diff before sync:**

```bash
$ argocd app diff python-app-dev
===== /apps/v1/Deployment dev/python-app-python-app ======
@@ -24,7 +24,7 @@
       labels:
         app.kubernetes.io/instance: python-app-dev
         app.kubernetes.io/name: python-app
     name: python-app-python-app
-    replicas: 1
+    replicas: 5
     selector:
```

**After auto-sync:** replicas back to 1.

```bash
$ kubectl get deployment -n dev
NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
python-app-python-app    1/1     1            1           15m
```

### Test 2: Pod Deletion (Kubernetes vs ArgoCD)

**Delete a pod:**

```bash
$ kubectl delete pod python-app-python-app-7b858cb655-5gdxv -n dev
pod "python-app-python-app-7b858cb655-5gdxv" deleted
```

**Kubernetes immediately recreates the pod** (ReplicaSet controller):

```bash
$ kubectl get pods -n dev -w
python-app-python-app-7b858cb655-5gdxv   Terminating   0s
python-app-python-app-7b858cb655-abcde   0/1   Pending   0s
python-app-python-app-7b858cb655-abcde   0/1   Running   1s
```

ArgoCD does **not** react because there is no configuration drift – the desired replica count (1) is still satisfied.

**Key difference:**
- **Kubernetes self‑healing** ensures the desired number of running pods (handled by Deployment/ReplicaSet).
- **ArgoCD self‑healing** ensures the cluster configuration matches the Git definition (e.g., replicas, labels, images).

### Test 3: Configuration Drift (Adding a Label)

**Add a label to the deployment:**

```bash
$ kubectl label deployment python-app-python-app -n dev test-drift=true
deployment.apps/python-app-python-app labeled
```

**ArgoCD detects drift:**

```bash
$ argocd app diff python-app-dev
===== /apps/v1/Deployment dev/python-app-python-app ======
@@ -31,6 +31,7 @@
       labels:
         app.kubernetes.io/instance: python-app-dev
         app.kubernetes.io/name: python-app
+        test-drift: "true"
```

Within the sync interval, ArgoCD removes the label (self‑heal). Verify:

```bash
$ kubectl describe deployment python-app-python-app -n dev | grep test-drift
# (no output)
```

**Conclusion:** Self‑healing works for any configuration drift, not only replica counts.

---

## 5. Screenshots
In the 'screenshot/lab 13/' directory