package note_test

import (
	"context"
	"testing"

	"estudos-golang/internal/domain/note"
	"estudos-golang/internal/infrastructure/memory"
	usecase "estudos-golang/internal/usecase/note"
)

func TestService_CreateAndGet(t *testing.T) {
	svc := usecase.NewService(memory.NewNoteRepository())
	ctx := context.Background()

	created, err := svc.Create(ctx, usecase.CreateInput{Title: "Test", Content: "Body"})
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.Get(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Test" {
		t.Fatalf("want Test got %s", got.Title)
	}
}

func TestService_DeleteNotFound(t *testing.T) {
	svc := usecase.NewService(memory.NewNoteRepository())
	err := svc.Delete(context.Background(), "missing")
	if err != note.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
