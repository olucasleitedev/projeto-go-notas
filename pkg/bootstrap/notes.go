package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	domain "estudos-golang/internal/domain/note"
	"estudos-golang/internal/infrastructure/memory"
	"estudos-golang/internal/infrastructure/postgres"
	"estudos-golang/pkg/config"
	"estudos-golang/pkg/database"
)

func NotesRepository(ctx context.Context) (domain.Repository, *sql.DB, error) {
	dsn := config.EnvOr("NOTES_DATABASE_URL", "")
	if dsn == "" {
		return memory.NewNoteRepository(), nil, nil
	}

	db, err := database.Open(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}

	ddl, err := database.MigrationSQL("notes.sql")
	if err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("read notes migration: %w", err)
	}
	if err := database.Migrate(ctx, db, ddl); err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	return postgres.NewNoteRepository(db), db, nil
}
