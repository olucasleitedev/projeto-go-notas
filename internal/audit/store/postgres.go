package store

import (
	"context"
	"database/sql"

	"estudos-golang/pkg/events"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Append(ctx context.Context, evt events.NoteEvent) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO audit_events (event_type, note_id, title, occurred_at)
		VALUES ($1, $2, $3, $4)
	`, evt.Type, evt.NoteID, evt.Title, evt.Timestamp)
	return err
}

func (s *PostgresStore) List(ctx context.Context) ([]events.NoteEvent, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT event_type, note_id, title, occurred_at
		FROM audit_events ORDER BY occurred_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]events.NoteEvent, 0)
	for rows.Next() {
		var evt events.NoteEvent
		if err := rows.Scan(&evt.Type, &evt.NoteID, &evt.Title, &evt.Timestamp); err != nil {
			return nil, err
		}
		items = append(items, evt)
	}
	return items, rows.Err()
}
