# LAB 9

## Architecture Overview

### Description of your deployment architecture
#### Descrition:
- A five‑node Minikube cluster (single node in typical setup) running the app-python Deployment.
- Deployment manages 5 replicas (scalable) of the devops-info-service container.
- A NodePort Service (devops-service) exposes the application externally on a high port (e.g., 30080).

- Incoming traffic:
minikube service devops-service → NodePort → Service → Pod → Flask app on port 5000.

#### Resource Allocation:
Each pod requests 128Mi memory / 100m CPU, with limits 256Mi / 200m.
This ensures fair resource usage and prevents resource starvation.

## Manifest Files

### Deployment
- Replicas: 3 (initial) – provides basic high availability.
- Strategy: RollingUpdate with maxSurge=1, maxUnavailable=0 – ensures zero downtime during updates.
- Resources: Requests set based on average load; limits to avoid noisy neighbors.
- Probes:
    - - Liveness: /health – checks if the app is alive.
    - - Readiness: /ready – ensures the app is ready to serve traffic.
- Ports: containerPort: 5000 – Flask default.

### Service (service.yaml)
- Type: NodePort – suitable for local Minikube access.
- Selector: app: app-python – matches pods.
- Ports: port: 80 (internal), targetPort: 5000 (container port).
- No explicit nodePort – Kubernetes assigns a random high port (e.g., 30000–32767).

## Deployment Evidence

```sh
$ kubectl get all
NAME                             READY   STATUS    RESTARTS   AGE
pod/app-python-9696f6cb7-9b79j   1/1     Running   0          88s
pod/app-python-9696f6cb7-cnvzc   1/1     Running   0          54s
pod/app-python-9696f6cb7-gmfts   1/1     Running   0          62s
pod/app-python-9696f6cb7-h87hb   1/1     Running   0          69s
pod/app-python-9696f6cb7-lhccx   1/1     Running   0          77s

NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
service/devops-service   NodePort    10.97.228.249   <none>        80:30371/TCP   13m
service/kubernetes       ClusterIP   10.96.0.1       <none>        443/TCP        101m

NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/app-python   5/5     5            5           51m

NAME                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/app-python-66c7c6f5b8   0         0         0       51m
replicaset.apps/app-python-9696f6cb7    5         5         5       88s

$ kubectl get pods,svc -o wide
NAME                             READY   STATUS    RESTARTS   AGE    IP            NODE       NOMINATED NODE   READINESS GATES
pod/app-python-9696f6cb7-9b79j   1/1     Running   0          2m6s   10.244.0.14   minikube   <none>           <none>
pod/app-python-9696f6cb7-cnvzc   1/1     Running   0          92s    10.244.0.19   minikube   <none>           <none>
pod/app-python-9696f6cb7-gmfts   1/1     Running   0          100s   10.244.0.18   minikube   <none>           <none>
pod/app-python-9696f6cb7-h87hb   1/1     Running   0          107s   10.244.0.17   minikube   <none>           <none>
pod/app-python-9696f6cb7-lhccx   1/1     Running   0          115s   10.244.0.16   minikube   <none>           <none>

NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE    SELECTOR
service/devops-service   NodePort    10.97.228.249   <none>        80:30371/TCP   14m    app=app-python
service/kubernetes       ClusterIP   10.96.0.1       <none>        443/TCP        102m   <none>

$ kubectl describe deployment app-python
Name:                   app-python
Namespace:              default
CreationTimestamp:      Thu, 26 Mar 2026 15:39:12 +0300
Labels:                 app=app-python
                        environment=production
Annotations:            deployment.kubernetes.io/revision: 2
Selector:               app=app-python
Replicas:               5 desired | 5 updated | 5 total | 5 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  0 max unavailable, 1 max surge
Pod Template:
  Labels:  app=app-python
  Containers:
   app-python:
    Image:      damirsadykov/devops-info-service:1.3.0
    Port:       5000/TCP (http)
    Host Port:  0/TCP (http)
    Limits:
      cpu:     200m
      memory:  256Mi
    Requests:
      cpu:      100m
      memory:   128Mi
    Liveness:   http-get http://:5000/health delay=10s timeout=1s period=5s #success=1 #failure=3
    Readiness:  http-get http://:5000/ready delay=5s timeout=1s period=3s #success=1 #failure=3
    Environment:
      PORT:        5000
      HOST:        0.0.0.0
      DEBUG:       false
    Mounts:        <none>
  Volumes:         <none>
  Node-Selectors:  <none>
  Tolerations:     <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  app-python-66c7c6f5b8 (0/0 replicas created)
NewReplicaSet:   app-python-9696f6cb7 (5/5 replicas created)
Events:
  Type    Reason             Age                  From                   Message
  ----    ------             ----                 ----                   -------
  Normal  ScalingReplicaSet  52m                  deployment-controller  Scaled up replica set app-python-66c7c6f5b8 from 0 to 3
  Normal  ScalingReplicaSet  2m26s                deployment-controller  Scaled up replica set app-python-66c7c6f5b8 from 3 to 5
  Normal  ScalingReplicaSet  2m26s                deployment-controller  Scaled up replica set app-python-9696f6cb7 from 0 to 1
  Normal  ScalingReplicaSet  2m15s                deployment-controller  Scaled down replica set app-python-66c7c6f5b8 from 5 to 4
  Normal  ScalingReplicaSet  2m15s                deployment-controller  Scaled up replica set app-python-9696f6cb7 from 1 to 2
  Normal  ScalingReplicaSet  2m7s                 deployment-controller  Scaled down replica set app-python-66c7c6f5b8 from 4 to 3
  Normal  ScalingReplicaSet  2m7s                 deployment-controller  Scaled up replica set app-python-9696f6cb7 from 2 to 3
  Normal  ScalingReplicaSet  2m                   deployment-controller  Scaled down replica set app-python-66c7c6f5b8 from 3 to 2
  Normal  ScalingReplicaSet  2m                   deployment-controller  Scaled up replica set app-python-9696f6cb7 from 3 to 4
  Normal  ScalingReplicaSet  112s                 deployment-controller  Scaled down replica set app-python-66c7c6f5b8 from 2 to 1
  Normal  ScalingReplicaSet  104s (x2 over 112s)  deployment-controller  (combined from similar events): Scaled down replica set app-python-66c7c6f5b8 from 1 to 0
```

## Operations Performed

### Commands used to deploy
```
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

kubectl get all
```
### Scaling demonstration output
```
$ kubectl scale deployment app-python --replicas=5
deployment.apps/app-python scaled

$ kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
app-python-9696f6cb7-9b79j   1/1     Running   0          5m42s
app-python-9696f6cb7-cnvzc   1/1     Running   0          5m8s
app-python-9696f6cb7-gmfts   1/1     Running   0          5m16s
app-python-9696f6cb7-h87hb   1/1     Running   0          5m23s
app-python-9696f6cb7-lhccx   1/1     Running   0          5m31s
```
### Rolling update demonstration output
```
$ kubectl set image deployment/app-python app-python=damirsadykov/devops-info-service:latest
deployment.apps/app-python image updated

$ kubectl rollout status deployment/app-python
Waiting for deployment "app-python" rollout to finish: 2 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 2 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 3 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 3 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 3 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 3 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 4 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 4 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 4 out of 5 new replicas have been updated...
Waiting for deployment "app-python" rollout to finish: 1 old replicas are pending termination...
Waiting for deployment "app-python" rollout to finish: 1 old replicas are pending termination...
Waiting for deployment "app-python" rollout to finish: 1 old replicas are pending termination...
deployment "app-python" successfully rolled out
```
### Service access method and verification
```
$ kubectl get svc devops-service

NAME             TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
devops-service   NodePort   10.97.228.249   <none>        80:30371/TCP   20m

$ minikube service devops-service
┌───────────┬────────────────┬─────────────┬───────────────────────────┐
│ NAMESPACE │      NAME      │ TARGET PORT │            URL            │
├───────────┼────────────────┼─────────────┼───────────────────────────┤
│ default   │ devops-service │ 80          │ http://192.168.49.2:30371 │
└───────────┴────────────────┴─────────────┴───────────────────────────┘
🎉  Opening service default/devops-service in default browser...
```

## Production Considerations

### Health Checks
- Liveness probe (/health): returns 200 if the app is alive.
- Readiness probe (/ready): returns 200 after app fully starts (e.g., after loading dependencies, establishing connections).
These ensure Kubernetes can restart unhealthy pods and route traffic only to ready ones.

### Resource Limits Rationale
- Requests: Guarantee minimum resources (128Mi / 100m) – enough for idle state.
- Limits: Prevent pods from consuming excessive resources (256Mi / 200m) – protects other workloads on the node.
- Values chosen based on typical memory/CPU usage of a Flask app with Prometheus metrics.

### Improvements for Production
- Use a LoadBalancer or Ingress instead of NodePort for external access.
- Enable Horizontal Pod Autoscaler (HPA) based on CPU/memory or custom metrics.
- ConfigMap / Secret for configuration (e.g., DEBUG, API keys).
- PodDisruptionBudget to ensure availability during voluntary disruptions.
- Network Policies to restrict ingress/egress traffic.
- Persistent storage if needed.
- Use a service mesh (e.g., Istio) for advanced traffic management.

### Monitoring and Observability Strategy
- Metrics: Already exposed at /metrics for Prometheus.
- Logging: JSON logs (as implemented) to stdout – easily collected by Fluentd/Elasticsearch or Loki.
- Dashboard: Grafana dashboards for Prometheus metrics.
- Tracing: Add OpenTelemetry for distributed tracing (e.g., Jaeger).
- Alerts: Set up alerts for high latency, pod restarts, or CPU/memory saturation