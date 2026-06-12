package note

import (
	"strings"
	"time"
)

// Note é a entidade do domínio — regras de negócio vivem aqui.
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(title, content string, now time.Time) (Note, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Note{}, ErrInvalidTitle
	}

	now = now.UTC()
	return Note{
		Title:     title,
		Content:   strings.TrimSpace(content),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (n *Note) Update(title, content string, now time.Time) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrInvalidTitle
	}

	n.Title = title
	n.Content = strings.TrimSpace(content)
	n.UpdatedAt = now.UTC()
	return nil
}
