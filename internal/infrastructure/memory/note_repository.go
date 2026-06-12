package memory

import (
	"context"
	"slices"
	"sync"

	domain "estudos-golang/internal/domain/note"
)

// NoteRepository é um "adapter" — implementação concreta do port Repository.
type NoteRepository struct {
	mu    sync.RWMutex
	notes map[string]domain.Note
}

func NewNoteRepository() *NoteRepository {
	return &NoteRepository{notes: make(map[string]domain.Note)}
}

func (r *NoteRepository) Save(_ context.Context, note domain.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notes[note.ID] = note
	return nil
}

func (r *NoteRepository) FindByID(_ context.Context, id string) (domain.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	note, ok := r.notes[id]
	if !ok {
		return domain.Note{}, domain.ErrNotFound
	}
	return note, nil
}

func (r *NoteRepository) FindAll(_ context.Context) ([]domain.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.Note, 0, len(r.notes))
	for _, n := range r.notes {
		items = append(items, n)
	}

	slices.SortFunc(items, func(a, b domain.Note) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})
	return items, nil
}

func (r *NoteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notes[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.notes, id)
	return nil
}
