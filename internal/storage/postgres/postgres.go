package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(p DBParams) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), p.URL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	if err = dbpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("could not ping postgres: %w", err)
	}

	return dbpool, nil
}
