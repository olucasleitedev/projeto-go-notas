package store

import (
	"context"
	"slices"
	"sync"

	"estudos-golang/pkg/events"
)

type MemoryStore struct {
	mu     sync.RWMutex
	events []events.NoteEvent
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{events: make([]events.NoteEvent, 0)}
}

func (s *MemoryStore) Append(_ context.Context, evt events.NoteEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, evt)
	return nil
}

func (s *MemoryStore) List(_ context.Context) ([]events.NoteEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]events.NoteEvent, len(s.events))
	copy(out, s.events)
	slices.SortFunc(out, func(a, b events.NoteEvent) int {
		return b.Timestamp.Compare(a.Timestamp)
	})
	return out, nil
}
