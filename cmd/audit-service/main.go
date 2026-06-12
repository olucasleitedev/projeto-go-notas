package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	audithandler "estudos-golang/internal/audit/handler"
	auditstore "estudos-golang/internal/audit/store"
	"estudos-golang/pkg/config"
	"estudos-golang/pkg/events"
	"estudos-golang/pkg/messaging"
)

func main() {
	store := auditstore.NewMemoryStore()
	h := audithandler.New(store)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if config.KafkaEnabled() {
		go func() {
			log.Printf("audit consumer starting (topic=%s)", events.TopicNoteEvents)
			err := messaging.Consume(ctx, config.KafkaBrokers(), "audit-service", events.TopicNoteEvents,
				func(ctx context.Context, _ []byte, value []byte) error {
					var evt events.NoteEvent
					if err := messaging.UnmarshalJSON(value, &evt); err != nil {
						return err
					}
					store.Append(evt)
					log.Printf("audit event: %s note_id=%s", evt.Type, evt.NoteID)
					return nil
				})
			if err != nil && ctx.Err() == nil {
				log.Printf("audit consumer stopped: %v", err)
			}
		}()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/audit/events", h.ListEvents)

	addr := config.EnvOr("ADDR", ":8082")
	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		log.Printf("audit-service listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	cancel()
	_ = server.Shutdown(context.Background())
}
