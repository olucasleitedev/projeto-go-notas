package note_test

import (
	"testing"
	"time"

	"estudos-golang/internal/domain/note"
)

func TestNew_ValidatesTitle(t *testing.T) {
	_, err := note.New("   ", "body", time.Now())
	if err != note.ErrInvalidTitle {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestUpdate_ValidatesTitle(t *testing.T) {
	n, err := note.New("ok", "body", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := n.Update("", "x", time.Now()); err != note.ErrInvalidTitle {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}
