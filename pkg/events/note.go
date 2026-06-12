package events

import "time"

const (
	TopicNoteEvents = "note.events"

	NoteCreated = "note.created"
	NoteUpdated = "note.updated"
	NoteDeleted = "note.deleted"
)

type NoteEvent struct {
	Type      string    `json:"type"`
	NoteID    string    `json:"note_id"`
	Title     string    `json:"title,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
