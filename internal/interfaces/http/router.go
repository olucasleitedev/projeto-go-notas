package http

import (
	"net/http"

	"estudos-golang/internal/interfaces/http/handler"
	"estudos-golang/internal/interfaces/http/middleware"
	usecase "estudos-golang/internal/usecase/note"
)

func NewNotesRouter(noteSvc *usecase.Service) http.Handler {
	mux := http.NewServeMux()
	notes := handler.NewNoteHandler(noteSvc)

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /api/notes", notes.List)
	mux.HandleFunc("POST /api/notes", notes.Create)
	mux.HandleFunc("GET /api/notes/{id}", notes.Get)
	mux.HandleFunc("PUT /api/notes/{id}", notes.Update)
	mux.HandleFunc("DELETE /api/notes/{id}", notes.Delete)

	return mux
}

// NewRouter mantém compatibilidade com o modo monólito (CORS incluso).
func NewRouter(noteSvc *usecase.Service) http.Handler {
	return middleware.CORS(NewNotesRouter(noteSvc))
}
