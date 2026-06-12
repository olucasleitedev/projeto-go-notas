package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"estudos-golang/internal/audit/store"
	"estudos-golang/pkg/events"
)

type Handler struct {
	store store.EventStore
	db    *sql.DB
}

func New(s store.EventStore) *Handler {
	return &Handler{store: s}
}

func NewWithDB(s store.EventStore, db *sql.DB) *Handler {
	return &Handler{store: s, db: db}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "audit"})
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready", "service": "audit"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready", "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready", "service": "audit"})
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list events"})
		return
	}
	if items == nil {
		items = []events.NoteEvent{}
	}
	writeJSON(w, http.StatusOK, items)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
