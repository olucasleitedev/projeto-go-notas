package handler

import (
	"encoding/json"
	"net/http"

	domain "estudos-golang/internal/domain/note"
	usecase "estudos-golang/internal/usecase/note"
)

// NoteHandler é um "driving adapter" — traduz HTTP em chamadas de caso de uso.
type NoteHandler struct {
	svc *usecase.Service
}

func NewNoteHandler(svc *usecase.Service) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in usecase.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "invalid json"})
		return
	}

	note, err := h.svc.Create(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, note)
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	notes, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	if notes == nil {
		notes = []domain.Note{}
	}
	writeJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) Get(w http.ResponseWriter, r *http.Request) {
	note, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	var in usecase.UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "invalid json"})
		return
	}

	note, err := h.svc.Update(r.Context(), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
