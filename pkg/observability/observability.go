package observability

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service struct {
	Name   string
	Logger *slog.Logger
}

func New(name string) *Service {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	return &Service{Name: name, Logger: logger}
}

type metrics struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func newMetrics(service string) *metrics {
	return &metrics{
		requests: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
			ConstLabels: prometheus.Labels{"service": service},
		}, []string{"method", "path", "status"}),
		duration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{"service": service},
		}, []string{"method", "path"}),
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (s *Service) Middleware(m *metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			path := r.Pattern
			if path == "" {
				path = r.URL.Path
			}
			m.requests.WithLabelValues(r.Method, path, strconv.Itoa(wrapped.status)).Inc()
			m.duration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())

			s.Logger.Info("request",
				"service", s.Name,
				"method", r.Method,
				"path", path,
				"status", wrapped.status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

func (s *Service) WrapHandler(m *metrics, handler http.Handler) http.Handler {
	return s.Middleware(m)(handler)
}

func MountMetrics(mux *http.ServeMux) {
	mux.Handle("GET /metrics", promhttp.Handler())
}

func NewMetrics(service string) *metrics {
	return newMetrics(service)
}
