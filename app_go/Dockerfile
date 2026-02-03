FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o myapp

FROM alpine:3.18
RUN adduser -D appuser
WORKDIR /app
COPY --from=builder /app/myapp .
USER appuser
EXPOSE 5000
CMD ["./myapp"]