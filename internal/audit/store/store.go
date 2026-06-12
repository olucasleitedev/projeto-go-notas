package store

import (
	"context"

	"estudos-golang/pkg/events"
)

type EventStore interface {
	Append(ctx context.Context, evt events.NoteEvent) error
	List(ctx context.Context) ([]events.NoteEvent, error)
}
