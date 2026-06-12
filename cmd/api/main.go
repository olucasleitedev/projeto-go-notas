package main

import (
	"log"
	"net/http"
	"os"

	apphttp "estudos-golang/internal/interfaces/http"
	"estudos-golang/internal/infrastructure/memory"
	usecase "estudos-golang/internal/usecase/note"
)

func main() {
	// Composition Root — único lugar que "liga" todas as camadas.
	repo := memory.NewNoteRepository()
	noteSvc := usecase.NewService(repo)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("API rodando em http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, apphttp.NewRouter(noteSvc)))
}
