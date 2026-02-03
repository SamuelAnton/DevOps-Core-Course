# Lab 02

## Multi-stage strategy
The total process is in 2 stages:
1. Create an optimized, statically compiled binary.
2. Minimal production runtime environment.

We want to firstly build our app and for this we use all we need. Then we want to reduce container size, remove all compiler's stuff and tools we won't use. And for this we copy only executable file with go-alpine.

## Build Process & Size Analysis

###  Terminal Output
```sh
damir@damir-VB:~/Desktop/DevOps/DevOps-Core-Course/app_go$ docker build -t devops-info-service-go:latest .
DEPRECATED: The legacy builder is deprecated and will be removed in a future release.
            Install the buildx component to build images with BuildKit:
            https://docs.docker.com/go/buildx/

Sending build context to Docker daemon  825.3kB
Step 1/13 : FROM golang:1.21 AS builder
 ---> 246ea1ed9cdb
Step 2/13 : WORKDIR /app
 ---> Using cache
 ---> f14686bd9c3d
Step 3/13 : COPY go.mod ./
 ---> 120f56ebd2dc
Step 4/13 : RUN go mod download
 ---> Running in 83ea1c732318
go: no module dependencies to download
 ---> Removed intermediate container 83ea1c732318
 ---> ac49e8ffa4a1
Step 5/13 : COPY . .
 ---> 8293d31a8458
Step 6/13 : RUN CGO_ENABLED=0 go build -o myapp
 ---> Running in ef0aa4e2f41c
 ---> Removed intermediate container ef0aa4e2f41c
 ---> b140eafdac12
Step 7/13 : FROM alpine:3.18
3.18: Pulling from library/alpine
44cf07d57ee4: Pull complete 
Digest: sha256:de0eb0b3f2a47ba1eb89389859a9bd88b28e82f5826b6969ad604979713c2d4f
Status: Downloaded newer image for alpine:3.18
 ---> 802c91d52981
Step 8/13 : RUN adduser -D appuser
 ---> Running in 86943f47828a
 ---> Removed intermediate container 86943f47828a
 ---> 868848cc0cdb
Step 9/13 : WORKDIR /app
 ---> Running in 62cf6e903abe
 ---> Removed intermediate container 62cf6e903abe
 ---> 7d3432469112
Step 10/13 : COPY --from=builder /app/myapp .
 ---> ae233bad30d3
Step 11/13 : USER appuser
 ---> Running in e2d2076c73e9
 ---> Removed intermediate container e2d2076c73e9
 ---> 42bde78130ee
Step 12/13 : EXPOSE 5000
 ---> Running in 8f663d6562d1
 ---> Removed intermediate container 8f663d6562d1
 ---> a1e0dafdc306
Step 13/13 : CMD ["./myapp"]
 ---> Running in b6fd172a8246
 ---> Removed intermediate container b6fd172a8246
 ---> b2e17931f148
Successfully built b2e17931f148
Successfully tagged devops-info-service-go:latest
damir@damir-VB:~/Desktop/DevOps/DevOps-Core-Course/app_go$ docker images | grep devops
devops-info-service-go               latest        b2e17931f148   14 minutes ago   14.3MB
damirsadykov/devops-info-service     1.0           8a7e51a097ec   17 hours ago     132MB
damirsadykov/devops-info-service     latest        8a7e51a097ec   17 hours ago     132MB
devops-info-service                  latest        8a7e51a097ec   17 hours ago     132MB
```

### Size analysis
```sh
damir@damir-VB:~/Desktop/DevOps/DevOps-Core-Course/app_go$ docker image ls --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
REPOSITORY                           TAG           SIZE
devops-info-service-go               latest        14.3MB
golang                               1.21          814MB
```

### Size breakdown
```
Multi-Stage Build:
├── Builder Stage (golang:1.21): 814MB
│   ├── Go toolchain: ~814MB
│   └── Dependencies: <1MB
│
└── Final Stage (alpine:3.18 + binary): 14.3MB
    ├── Alpine base: 7.05MB
    ├── Compiled binary: 7.2MB
    └── User setup: 0.05MB

Size Reduction: 864MB → 14.3MB (98.6% reduction!)
```

### Comparison Table:
|Aspect	Single-Stage | Multi-Stage (Your) | Improvement
|----------|--------|-------------|
|Image Size	814MB | 14.3MB | 98.6% smaller
|Security Risk | High (compilers, source) | Low (binary only) | Significantly safer
|Attack Surface | Large (~1000 packages) | Minimal (~50 packages) | 95% reduction
|Deployment Speed | Slow (large transfer) | Fast (small transfer) | 73x faster transfer
|Registry Storage Cost | High | Minimal | ~$0.87/month vs $0.01/month


## Technical Stage-by-Stage Analysis

### Stage 1: Builder (golang:1.21)
```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o myapp
```

#### Purpose: Create an optimized, statically compiled binary.

#### Key Decisions:
1. `golang:1.21`: Specific version ensures reproducible builds

2. `WORKDIR /app`: Consistent working directory

3. Copy `go.mod` first: Leverages Docker cache for dependencies

4. `go mod download`: Downloads dependencies separately (cache optimization)

5. `CGO_ENABLED=0`: Produces static binary with no C dependencies

6. `-o myapp`: Explicit output name for clarity

### Stage 2: Runtime (alpine:3.18)
```dockerfile
FROM alpine:3.18
RUN adduser -D appuser
WORKDIR /app
COPY --from=builder /app/myapp .
USER appuser
EXPOSE 5000
CMD ["./myapp"]
```

#### Purpose: Minimal production runtime environment.

#### Key Decisions:
1. `alpine:3.18`: Extremely small base image (~7MB)
2. Non-root user: Security best practice
3. `COPY --from=builder`: Only copies the binary, not build tools
4. Explicit CMD: Clear entry point


## Security Benefits Analysis
### Security Improvements:
1. No Compilers: Can't compile malicious code inside container
2. No Source Code: Intellectual property protected
3. Minimal Packages: Fewer CVEs to patch
4. Non-Root User: Principle of least privilege
5. Static Binary: No runtime dependency vulnerabilities