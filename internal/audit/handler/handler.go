package handler

import (
	"encoding/json"
	"net/http"

	"estudos-golang/internal/audit/store"
)

type Handler struct {
	store *store.MemoryStore
}

func New(store *store.MemoryStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "audit"})
}

func (h *Handler) ListEvents(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.store.List())
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
