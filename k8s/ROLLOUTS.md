# Lab 14

## 1. Argo Rollouts Setup

### Installation Verification

Argo Rollouts controller and kubectl plugin were installed in the `argo-rollouts` namespace.

```bash
$ kubectl get pods -n argo-rollouts
NAME                                      READY   STATUS    RESTARTS   AGE
argo-rollouts-dashboard-755bbc64c-jq8hr   1/1     Running   0          81m
argo-rollouts-ffb479f4-2jdzs              1/1     Running   0          105m
argo-rollouts-ffb479f4-rhj6q              1/1     Running   0          105m

$ kubectl argo rollouts version
kubectl-argo-rollouts: v1.9.0+838d4e7
  BuildDate: 2026-03-20T21:08:11Z
  GitCommit: 838d4e792be666ec11bd0c80331e0c5511b5010e
  GitTreeState: clean
  GoVersion: go1.24.13
  Compiler: gc
  Platform: linux/amd64
```

Rollout CRD installed:

```bash
$ kubectl get crd | grep rollouts
rollouts.argoproj.io                   2026-04-30T11:13:39Z
```

### Dashboard Access

Dashboard was deployed and accessed via port-forward:

```bash
$ kubectl port-forward svc/argo-rollouts-dashboard -n argo-rollouts 3100:3100
Forwarding from 127.0.0.1:3100 -> 3100
```

Screenshot: screenshots/lab 14/Dashboard access.png

---

## 2. Canary Deployment

### Strategy Configuration

The Helm chart `templates/rollout.yaml` (after 3rd task `templates/rollout canary.yaml.txt`) contains the following canary steps:

```yaml
strategy:
  canary:
    steps:
      - setWeight: 20
      - pause: {}
      - setWeight: 40
      - pause: { duration: 30s }
      - setWeight: 60
      - pause: { duration: 30s }
      - setWeight: 80
      - pause: { duration: 30s }
      - setWeight: 100
```

- **20% weight** → manual pause – operator must verify before proceeding.
- Subsequent weights (40%, 60%, 80%) automatically continue after 30 seconds.
- Final step switches 100% to new version.

### Step-by-Step Rollout Progression - screenshots inside `screensots/` folder

### Promotion and Abort Demonstration

#### Promote to next step

```bash
$ kubectl argo rollouts promote my-python-app
rollout "my-python-app" promoted
```

#### Abort during canary (e.g., at 20% weight)

```bash
$ kubectl argo rollouts abort my-python-app
rollout "my-python-app" aborted
```

Abort instantly reverts all traffic to the stable version. Dashboard shows `Aborted` status.

To retry after abort:

```bash
$ kubectl argo rollouts retry my-python-app
```

---

## 3. Blue-Green Deployment

### Strategy Configuration

Blue-green strategy defined in the same `rollout.yaml`:

```yaml
strategy:
  blueGreen:
    activeService: my-python-app-python-app
    previewService: my-python-app-python-app-preview
    autoPromotionEnabled: false
```

- `activeService` → production traffic.
- `previewService` → points only to the new (green) version for testing.
- `autoPromotionEnabled: false` – must manually promote after verification.

### Preview vs Active Service

After updating the image tag, the rollout creates a new ReplicaSet (green) and the preview service targets it. The active service still targets the old (blue) version.

```bash
$ kubectl get svc
NAME                               TYPE        CLUSTER-IP      PORT(S)
my-python-app-python-app           NodePort    10.99.123.45    80:30080/TCP
my-python-app-python-app-preview   NodePort    10.99.123.46    80:30081/TCP
```

Access both services for comparison:

```bash
# Active (blue) version
$ curl http://$(minikube ip):30080/version
{"version":"1.0.0"}

# Preview (green) version
$ curl http://$(minikube ip):30081/version
{"version":"1.1.0"}
```

### Promotion Process

After testing the green version via preview service, promote it to active:

```bash
$ kubectl argo rollouts promote my-python-app
rollout "my-python-app" promoted
```

The active service instantly switches to the green ReplicaSet. No gradual traffic shift – all users see the new version immediately.

**Instant rollback** (switch back to blue):

```bash
$ kubectl argo rollouts undo my-python-app-python-app
```

---

## 4. Strategy Comparison

| Aspect               | Canary                                       | Blue-Green                                    |
|----------------------|----------------------------------------------|-----------------------------------------------|
| **Traffic switching** | Gradual (percentage weights)                | All-or-nothing (instant)                      |
| **Resource usage**   | Incremental – only small % of new version   | Requires **2x** resources (full blue + green) |
| **Rollback speed**   | Gradual (or instant if abort)               | Instant (switch service selector)             |
| **Testing before production** | Only partial real traffic          | Full preview environment with isolated access |
| **Risk**             | Lower – issues affect only a fraction       | Higher – total switch, but can rollback fast  |
| **Ideal for**        | – Low-risk updates<br>– A/B testing<br>– Gradual confidence building | – High-risk updates<br>– Require full validation before production<br>– Regulatory/compliance needs |

### When to Use Which

- **Canary**  
  - You want to test with real user traffic gradually.  
  - You have few resources and cannot double capacity.  
  - Your application can run two different versions simultaneously without conflict (e.g., stateless, backward compatible).

- **Blue-Green**  
  - You need a **full isolated environment** to test the new version.  
  - You require **instant rollback** (with zero downtime).  
  - Resources are not a constraint (you can run two full environments during deployment).  
  - Database changes must be backward compatible or handled separately (e.g., using feature flags).

### Recommendation

- For **internal tools / low‑risk services** → Canary is efficient and safe.  
- For **customer‑facing, mission‑critical applications** → Blue‑green provides the highest assurance, as every change is validated in a production‑like preview before any user sees it.  
- If resource limits are strict → Canary is more economical.

---

## 5. CLI Commands Reference

| Command | Description |
|---------|-------------|
| `kubectl argo rollouts version` | Show plugin and controller version. |
| `kubectl argo rollouts get rollout <name>` | Display current rollout status, weights, and steps. |
| `kubectl argo rollouts get rollout <name> --watch` | Live watch of rollout progression. |
| `kubectl argo rollouts history <name>` | Show rollout history and revisions. |
| `kubectl argo rollouts promote <name>` | Manually advance to the next canary step or promote green to active in blue‑green. |
| `kubectl argo rollouts abort <name>` | Abort an ongoing rollout, reverting traffic to stable version. |
| `kubectl argo rollouts retry <name>` | Retry a previously aborted rollout. |
| `kubectl argo rollouts undo <name>` | Rollback to a previous revision (instant in blue‑green). |
| `kubectl argo rollouts set image <name> <container>=<image>` | Quickly update the rollout image. |
| `kubectl describe rollout <name>` | Detailed event and condition information. |
| `kubectl get rollouts --all-namespaces` | List all Rollouts across namespaces. |

### Monitoring and Troubleshooting

- **Dashboard**: Provides a clean visual of canary steps, weights, and promotion history.  
- **Logs**: Check rollout controller logs if something behaves unexpectedly:
  ```bash
  kubectl logs -n argo-rollouts deployment/argo-rollouts
  ```
- **Events**:
  ```bash
  kubectl describe rollout <name> | grep -A 20 Events
  ```
- **ReplicaSets**: Compare `rollouts-pod-template-hash` labels to see which ReplicaSet is active/preview.


## Bonus: Automated Analysis

### AnalysisTemplate Definition

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AnalysisTemplate
```

### Integration with Canary
The analysis step runs after the 20% manual promotion. It queries /health three times (5s interval). If any measurement fails (status != "healthy"), the rollout aborts.

### Success Criteria
- Success: status == "healthy" for all 3 measurements.
- Failure limit: 1 failure → abort.

### Screenshots - inside `screenshots/` folder