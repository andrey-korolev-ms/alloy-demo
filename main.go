package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "demo_requests_total",
	Help: "Total /hello requests",
}, []string{"path"})
var requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "demo_requests_duration_seconds",
	Help:    "Latency of /hello",
	Buckets: prometheus.DefBuckets,
}, []string{"path"})

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Ошибка .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requests.WithLabelValues("/hello").Inc()
		w.Write([]byte("hello"))
		requestDuration.WithLabelValues("/hello").Observe(time.Since(start).Seconds())
	})
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Starting server on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen and serve error: %v", err)
		}
	}()

	// Обработка сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
