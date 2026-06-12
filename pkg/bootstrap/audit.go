package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	"estudos-golang/internal/audit/store"
	"estudos-golang/pkg/config"
	"estudos-golang/pkg/database"
)

func AuditStore(ctx context.Context) (store.EventStore, *sql.DB, error) {
	dsn := config.EnvOr("AUDIT_DATABASE_URL", "")
	if dsn == "" {
		return store.NewMemoryStore(), nil, nil
	}

	db, err := database.Open(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}

	ddl, err := database.MigrationSQL("audit.sql")
	if err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("read audit migration: %w", err)
	}
	if err := database.Migrate(ctx, db, ddl); err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	return store.NewPostgresStore(db), db, nil
}
