package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	apphttp "estudos-golang/internal/interfaces/http"
	"estudos-golang/internal/infrastructure/memory"
	usecase "estudos-golang/internal/usecase/note"
	"estudos-golang/pkg/config"
	"estudos-golang/pkg/messaging"
)

func main() {
	repo := memory.NewNoteRepository()

	var publisher messaging.Publisher = messaging.NoopPublisher{}
	if config.KafkaEnabled() {
		pub := messaging.NewKafkaPublisher(config.KafkaBrokers())
		publisher = pub
		defer pub.Close()
		log.Printf("kafka publisher enabled: %v", config.KafkaBrokers())
	}

	noteSvc := usecase.NewServiceWithEvents(repo, publisher)

	addr := config.EnvOr("ADDR", ":8081")
	server := &http.Server{
		Addr:    addr,
		Handler: apphttp.NewNotesRouter(noteSvc),
	}

	go func() {
		log.Printf("notes-service listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	_ = server.Shutdown(context.Background())
}
