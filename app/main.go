package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	appName     = getEnv("APP_NAME", "gitops-sample-service")
	environment = getEnv("APP_ENV", "development")
	version     = getEnv("APP_VERSION", "1.0.0")
	port        = getEnv("PORT", "8080")

	// Atomic counters for simple Prometheus-compatible metrics
	requestCounter uint64
	healthyState   int32 = 1
	readyState     int32 = 1
)

type StatusResponse struct {
	Status      string    `json:"status"`
	App         string    `json:"app"`
	Environment string    `json:"environment"`
	Version     string    `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	UptimeSec   float64   `json:"uptime_seconds"`
}

var startTime = time.Now()

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&requestCounter, 1)
	w.Header().Set("Content-Type", "application/json")
	resp := StatusResponse{
		Status:      "operational",
		App:         appName,
		Environment: environment,
		Version:     version,
		Timestamp:   time.Now().UTC(),
		UptimeSec:   time.Since(startTime).Seconds(),
	}
	json.NewEncoder(w).Encode(resp)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthyState) == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("UNHEALTHY\n"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&readyState) == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY\n"))
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("NOT READY\n"))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	totalReqs := atomic.LoadUint64(&requestCounter)
	uptime := time.Since(startTime).Seconds()

	metrics := fmt.Sprintf(
		"# HELP http_requests_total Total number of HTTP requests\n"+
			"# TYPE http_requests_total counter\n"+
			"http_requests_total{app=\"%s\",env=\"%s\"} %d\n"+
			"# HELP app_uptime_seconds Process uptime in seconds\n"+
			"# TYPE app_uptime_seconds gauge\n"+
			"app_uptime_seconds{app=\"%s\",env=\"%s\"} %.2f\n"+
			"# HELP app_build_info Build and version details\n"+
			"# TYPE app_build_info gauge\n"+
			"app_build_info{app=\"%s\",version=\"%s\",env=\"%s\"} 1\n",
		appName, environment, totalReqs,
		appName, environment, uptime,
		appName, version, environment,
	)
	w.Write([]byte(metrics))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/ready", readyHandler)
	mux.HandleFunc("/metrics", metricsHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Trap SIGTERM and SIGINT for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting %s (%s) on port :%s (env: %s)", appName, version, port, environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sig := <-stopChan
	log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)

	// Mark not ready first so ingress / kube-proxy stops sending traffic
	atomic.StoreInt32(&readyState, 0)
	time.Sleep(2 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown error: %v", err)
	} else {
		log.Println("Server terminated cleanly.")
	}
}
