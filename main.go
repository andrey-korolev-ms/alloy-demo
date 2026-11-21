package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "demo_request_total",
	Help: "Total /hello requests",
}, []string{"path"})

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		requests.WithLabelValues("/hello").Inc()
		w.Write([]byte("hello"))
	})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
