package note

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	domain "estudos-golang/internal/domain/note"
)

// Service é a camada de aplicação (casos de uso).
// Orquestra o domínio e chama o repositório — sem saber de HTTP ou banco.
type Service struct {
	repo domain.Repository
	now  func() time.Time
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
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
	return n, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
