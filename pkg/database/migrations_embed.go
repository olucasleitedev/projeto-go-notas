package database

import "embed"

//go:embed migrations/*.sql
var migrationsFS embed.FS

func MigrationSQL(name string) (string, error) {
	b, err := migrationsFS.ReadFile("migrations/" + name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
