package note

import "context"

// Repository é um "port" — contrato que a infraestrutura implementa.
// O domínio não sabe se os dados estão em memória, Postgres ou arquivo.
type Repository interface {
	Save(ctx context.Context, note Note) error
	FindByID(ctx context.Context, id string) (Note, error)
	FindAll(ctx context.Context) ([]Note, error)
	Delete(ctx context.Context, id string) error
}
