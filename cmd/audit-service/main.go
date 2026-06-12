package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	audithandler "estudos-golang/internal/audit/handler"
	"estudos-golang/pkg/bootstrap"
	"estudos-golang/pkg/config"
	"estudos-golang/pkg/events"
	"estudos-golang/pkg/messaging"
	"estudos-golang/pkg/observability"
)

func main() {
	obs := observability.New("audit-service")
	metrics := observability.NewMetrics("audit-service")

	ctx := context.Background()
	eventStore, db, err := bootstrap.AuditStore(ctx)
	if err != nil {
		obs.Logger.Error("database init failed", "err", err)
		os.Exit(1)
	}
	if db != nil {
		defer db.Close()
		obs.Logger.Info("postgres connected")
	}

	var h *audithandler.Handler
	if db != nil {
		h = audithandler.NewWithDB(eventStore, db)
	} else {
		h = audithandler.New(eventStore)
	}

	consumerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if config.KafkaEnabled() {
		go func() {
			obs.Logger.Info("kafka consumer starting", "topic", events.TopicNoteEvents)
			err := messaging.Consume(consumerCtx, config.KafkaBrokers(), "audit-service", events.TopicNoteEvents,
				func(ctx context.Context, _ []byte, value []byte) error {
					var evt events.NoteEvent
					if err := messaging.UnmarshalJSON(value, &evt); err != nil {
						return err
					}
					if err := eventStore.Append(ctx, evt); err != nil {
						return err
					}
					obs.Logger.Info("event stored", "type", evt.Type, "note_id", evt.NoteID)
					return nil
				})
			if err != nil && consumerCtx.Err() == nil {
				obs.Logger.Error("consumer stopped", "err", err)
			}
		}()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /health/ready", h.Ready)
	observability.MountMetrics(mux)
	mux.HandleFunc("GET /api/audit/events", h.ListEvents)

	addr := config.EnvOr("ADDR", ":8082")
	server := &http.Server{Addr: addr, Handler: obs.WrapHandler(metrics, mux)}

	go func() {
		obs.Logger.Info("listening", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			obs.Logger.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	cancel()
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	_ = server.Shutdown(shutdownCtx)
	obs.Logger.Info("shutdown complete")
}
