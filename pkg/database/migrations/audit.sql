CREATE TABLE IF NOT EXISTS audit_events (
    id         BIGSERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    note_id    TEXT NOT NULL,
    title      TEXT NOT NULL DEFAULT '',
    occurred_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_occurred_at ON audit_events (occurred_at DESC);
