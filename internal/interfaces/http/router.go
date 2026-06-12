package http

import (
	"database/sql"
	"net/http"

	"estudos-golang/internal/interfaces/http/handler"
	"estudos-golang/internal/interfaces/http/middleware"
	usecase "estudos-golang/internal/usecase/note"
	"estudos-golang/pkg/observability"
)

func NewNotesRouter(noteSvc *usecase.Service, db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	notes := handler.NewNoteHandler(noteSvc)
	health := handler.NewHealthHandler(db)

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /health/ready", health.Ready)
	observability.MountMetrics(mux)
	mux.HandleFunc("GET /api/notes", notes.List)
	mux.HandleFunc("POST /api/notes", notes.Create)
	mux.HandleFunc("GET /api/notes/{id}", notes.Get)
	mux.HandleFunc("PUT /api/notes/{id}", notes.Update)
	mux.HandleFunc("DELETE /api/notes/{id}", notes.Delete)

	return mux
}

// NewRouter mantém compatibilidade com o modo monólito (CORS incluso).
func NewRouter(noteSvc *usecase.Service) http.Handler {
	return middleware.CORS(NewNotesRouter(noteSvc, nil))
}
