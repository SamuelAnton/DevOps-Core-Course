package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

// Structs for the response
type ServiceInfo struct {
	Service   Service    `json:"service"`
	System    System     `json:"system"`
	Runtime   Runtime    `json:"runtime"`
	Request   Request    `json:"request"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Service struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Framework   string `json:"framework"`
}

type System struct {
	Hostname        string `json:"hostname"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	Architecture    string `json:"architecture"`
	CPUCount        int    `json:"cpu_count"`
	GoVersion       string `json:"go_version"`
}

type Runtime struct {
	UptimeSeconds int    `json:"uptime_seconds"`
	UptimeHuman   string `json:"uptime_human"`
	CurrentTime   string `json:"current_time"`
	Timezone      string `json:"timezone"`
}

type Request struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
	Method    string `json:"method"`
	Path      string `json:"path"`
}

type Endpoint struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type HealthResponse struct {
	Status        string `json:"status"`
	Timestamp     string `json:"timestamp"`
	UptimeSeconds int    `json:"uptime_seconds"`
}

// Application start time
var startTime = time.Now()

// Helper function to get client IP
func getClientIP(r *http.Request) string {
	// Try to get IP from X-Forwarded-For header
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		// Fall back to RemoteAddr
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return r.RemoteAddr
		}
		return ip
	}
	return ip
}

// Helper function to format uptime
func getUptime() (int, string) {
	uptime := time.Since(startTime)
	seconds := int(uptime.Seconds())

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	if hours > 0 {
		return seconds, fmt.Sprintf("%d hour%s, %d minute%s",
			hours, plural(hours),
			minutes, plural(minutes))
	} else if minutes > 0 {
		return seconds, fmt.Sprintf("%d minute%s",
			minutes, plural(minutes))
	}
	return seconds, fmt.Sprintf("%d second%s",
		remainingSeconds, plural(remainingSeconds))
}

// Helper function for pluralization
func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

// Main handler
func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Get uptime
	uptimeSeconds, uptimeHuman := getUptime()

	// Get hostname
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	// Create response
	info := ServiceInfo{
		Service: Service{
			Name:        "devops-info-service",
			Version:     "1.0.0",
			Description: "DevOps course info service",
			Framework:   "Go",
		},
		System: System{
			Hostname:        hostname,
			Platform:        runtime.GOOS,
			PlatformVersion: runtime.Version(), // Go doesn't have direct OS version, using Go version
			Architecture:    runtime.GOARCH,
			CPUCount:        runtime.NumCPU(),
			GoVersion:       runtime.Version(),
		},
		Runtime: Runtime{
			UptimeSeconds: uptimeSeconds,
			UptimeHuman:   uptimeHuman,
			CurrentTime:   time.Now().UTC().Format(time.RFC3339Nano),
			Timezone:      "UTC",
		},
		Request: Request{
			ClientIP:  getClientIP(r),
			UserAgent: r.UserAgent(),
			Method:    r.Method,
			Path:      r.URL.Path,
		},
		Endpoints: []Endpoint{
			{Path: "/", Method: "GET", Description: "Service information"},
			{Path: "/health", Method: "GET", Description: "Health check"},
		},
	}

	// Set headers and encode JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)
}

// Health handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	uptimeSeconds, _ := getUptime()

	response := HealthResponse{
		Status:        "healthy",
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
		UptimeSeconds: uptimeSeconds,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Get configuration from environment variables
	host := getEnv("HOST", "0.0.0.0")
	port := getEnv("PORT", "5000")
	debug := getEnv("DEBUG", "false")

	// Set up routes
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/health", healthHandler)

	// Log startup information
	log.Printf("Starting DevOps Info Service (Go)...")
	log.Printf("Host: %s, Port: %s, Debug: %s", host, port, debug)
	log.Printf("Go version: %s", runtime.Version())
	log.Printf("Listening on http://%s:%s", host, port)

	// Start server
	addr := fmt.Sprintf("%s:%s", host, port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
