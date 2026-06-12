package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"estudos-golang/internal/interfaces/http/middleware"
	"estudos-golang/pkg/config"
	"estudos-golang/pkg/observability"
)

func main() {
	obs := observability.New("gateway")
	metrics := observability.NewMetrics("gateway")

	notesURL := mustParseURL(config.EnvOr("NOTES_SERVICE_URL", "http://localhost:8081"))
	auditURL := mustParseURL(config.EnvOr("AUDIT_SERVICE_URL", "http://localhost:8082"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","service":"gateway"}`))
	})
	observability.MountMetrics(mux)
	mux.Handle("/api/notes", proxy(notesURL, obs))
	mux.Handle("/api/notes/", proxy(notesURL, obs))
	mux.Handle("/api/audit/", proxy(auditURL, obs))

	addr := config.EnvOr("ADDR", ":8080")
	handler := middleware.CORS(obs.WrapHandler(metrics, mux))
	server := &http.Server{Addr: addr, Handler: handler}

	go func() {
		obs.Logger.Info("listening", "addr", addr, "notes", notesURL.String(), "audit", auditURL.String())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			obs.Logger.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
	obs.Logger.Info("shutdown complete")
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic("invalid url: " + raw)
	}
	return u
}

func proxy(target *url.URL, obs *observability.Service) http.Handler {
	rp := httputil.NewSingleHostReverseProxy(target)
	original := rp.Director
	rp.Director = func(req *http.Request) {
		original(req)
		req.Host = target.Host
	}
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		obs.Logger.Error("proxy error", "method", r.Method, "path", r.URL.Path, "err", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"upstream unavailable"}`, http.StatusBadGateway)
	}
	return rp
}
