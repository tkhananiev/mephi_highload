package metrics

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
)

var (
    TotalRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    RequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(TotalRequests)
    prometheus.MustRegister(RequestDuration)
}

type statusWriter struct {
    http.ResponseWriter
    status int
}

func (w *statusWriter) WriteHeader(code int) {
    w.status = code
    w.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sw := &statusWriter{ResponseWriter: w, status: 200}
        start := time.Now()
        next.ServeHTTP(sw, r)
        endpoint := r.URL.Path
        TotalRequests.WithLabelValues(r.Method, endpoint, http.StatusText(sw.status)).Inc()
        RequestDuration.WithLabelValues(r.Method, endpoint).Observe(time.Since(start).Seconds())
    })
}
