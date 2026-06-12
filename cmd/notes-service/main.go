package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apphttp "estudos-golang/internal/interfaces/http"
	usecase "estudos-golang/internal/usecase/note"
	"estudos-golang/pkg/bootstrap"
	"estudos-golang/pkg/config"
	"estudos-golang/pkg/messaging"
	"estudos-golang/pkg/observability"
)

func main() {
	obs := observability.New("notes-service")
	metrics := observability.NewMetrics("notes-service")

	ctx := context.Background()
	repo, db, err := bootstrap.NotesRepository(ctx)
	if err != nil {
		obs.Logger.Error("database init failed", "err", err)
		os.Exit(1)
	}
	if db != nil {
		defer db.Close()
		obs.Logger.Info("postgres connected")
	}

	var publisher messaging.Publisher = messaging.NoopPublisher{}
	if config.KafkaEnabled() {
		pub := messaging.NewKafkaPublisher(config.KafkaBrokers())
		publisher = pub
		defer pub.Close()
		obs.Logger.Info("kafka publisher enabled", "brokers", config.KafkaBrokers())
	}

	noteSvc := usecase.NewServiceWithEvents(repo, publisher)
	handler := obs.WrapHandler(metrics, apphttp.NewNotesRouter(noteSvc, db))

	addr := config.EnvOr("ADDR", ":8081")
	server := &http.Server{Addr: addr, Handler: handler}

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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
	obs.Logger.Info("shutdown complete")
}
