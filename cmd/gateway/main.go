package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"estudos-golang/internal/interfaces/http/middleware"
	"estudos-golang/pkg/config"
)

// API Gateway — ponto único de entrada para o frontend.

func main() {
	notesURL := mustParseURL(config.EnvOr("NOTES_SERVICE_URL", "http://localhost:8081"))
	auditURL := mustParseURL(config.EnvOr("AUDIT_SERVICE_URL", "http://localhost:8082"))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","service":"gateway"}`))
	})

	mux.Handle("/api/notes", proxy(notesURL))
	mux.Handle("/api/notes/", proxy(notesURL))
	mux.Handle("/api/audit/", proxy(auditURL))

	addr := config.EnvOr("ADDR", ":8080")
	server := &http.Server{
		Addr:    addr,
		Handler: middleware.CORS(mux),
	}

	go func() {
		log.Printf("gateway listening on %s", addr)
		log.Printf("proxy notes -> %s", notesURL)
		log.Printf("proxy audit -> %s", auditURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	_ = server.Shutdown(context.Background())
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("invalid url %q: %v", raw, err)
	}
	return u
}

func proxy(target *url.URL) http.Handler {
	rp := httputil.NewSingleHostReverseProxy(target)
	original := rp.Director
	rp.Director = func(req *http.Request) {
		original(req)
		req.Host = target.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	}
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error %s %s: %v", r.Method, r.URL.Path, err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"upstream unavailable"}`, http.StatusBadGateway)
	}
	return rp
}
