package postgres

import (
	"context"
	"database/sql"
	"errors"

	domain "estudos-golang/internal/domain/note"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Save(ctx context.Context, note domain.Note) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO notes (id, title, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			updated_at = EXCLUDED.updated_at
	`, note.ID, note.Title, note.Content, note.CreatedAt, note.UpdatedAt)
	return err
}

func (r *NoteRepository) FindByID(ctx context.Context, id string) (domain.Note, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes WHERE id = $1
	`, id)

	var n domain.Note
	err := row.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Note{}, domain.ErrNotFound
	}
	return n, err
}

func (r *NoteRepository) FindAll(ctx context.Context) ([]domain.Note, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Note, 0)
	for rows.Next() {
		var n domain.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, n)
	}
	return items, rows.Err()
}

func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
