//go:build integration

package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	domain "estudos-golang/internal/domain/note"
	"estudos-golang/internal/infrastructure/postgres"
	"estudos-golang/pkg/database"
)

func TestPostgresNoteRepository(t *testing.T) {
	dsn := os.Getenv("NOTES_DATABASE_URL")
	if dsn == "" {
		t.Skip("NOTES_DATABASE_URL not set")
	}

	ctx := context.Background()
	db, err := database.Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ddl, err := database.MigrationSQL("notes.sql")
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(ctx, db, ddl); err != nil {
		t.Fatal(err)
	}

	repo := postgres.NewNoteRepository(db)
	now := time.Now().UTC()
	n, err := domain.New("Postgres note", "content", now)
	if err != nil {
		t.Fatal(err)
	}
	n.ID = "integration-test-id"

	if err := repo.Save(ctx, n); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = repo.Delete(ctx, n.ID) }()

	got, err := repo.FindByID(ctx, n.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != n.Title {
		t.Fatalf("want %s got %s", n.Title, got.Title)
	}
}
