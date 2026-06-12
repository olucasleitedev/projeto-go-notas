package note

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	domain "estudos-golang/internal/domain/note"
	"estudos-golang/pkg/events"
	"estudos-golang/pkg/messaging"
)

type Service struct {
	repo      domain.Repository
	publisher messaging.Publisher
	now       func() time.Time
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func NewServiceWithEvents(repo domain.Repository, publisher messaging.Publisher) *Service {
	return &Service{repo: repo, publisher: publisher, now: time.Now}
}

type CreateInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *Service) Create(ctx context.Context, in CreateInput) (domain.Note, error) {
	n, err := domain.New(in.Title, in.Content, s.now())
	if err != nil {
		return domain.Note{}, err
	}

	n.ID = newID()
	if err := s.repo.Save(ctx, n); err != nil {
		return domain.Note{}, err
	}

	s.publish(ctx, events.NoteEvent{
		Type:      events.NoteCreated,
		NoteID:    n.ID,
		Title:     n.Title,
		Timestamp: s.now().UTC(),
	})
	return n, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Note, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (domain.Note, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, in UpdateInput) (domain.Note, error) {
	n, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Note{}, err
	}

	if err := n.Update(in.Title, in.Content, s.now()); err != nil {
		return domain.Note{}, err
	}

	if err := s.repo.Save(ctx, n); err != nil {
		return domain.Note{}, err
	}

	s.publish(ctx, events.NoteEvent{
		Type:      events.NoteUpdated,
		NoteID:    n.ID,
		Title:     n.Title,
		Timestamp: s.now().UTC(),
	})
	return n, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	n, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.publish(ctx, events.NoteEvent{
		Type:      events.NoteDeleted,
		NoteID:    n.ID,
		Title:     n.Title,
		Timestamp: s.now().UTC(),
	})
	return nil
}

func (s *Service) publish(ctx context.Context, evt events.NoteEvent) {
	if s.publisher == nil {
		return
	}
	_ = messaging.PublishNoteEvent(ctx, s.publisher, evt)
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
