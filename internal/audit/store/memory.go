package store

import (
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

func (s *MemoryStore) Append(evt events.NoteEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, evt)
}

func (s *MemoryStore) List() []events.NoteEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]events.NoteEvent, len(s.events))
	copy(out, s.events)
	slices.SortFunc(out, func(a, b events.NoteEvent) int {
		return b.Timestamp.Compare(a.Timestamp)
	})
	return out
}
