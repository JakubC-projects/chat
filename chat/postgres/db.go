package postgres

import (
	"context"

	"github.com/JakubC-projects/chat/chat/postgres/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	connPool *pgxpool.Pool
	queries  *sqlc.Queries
}

func NewDb(connStr string) *DB {
	connPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	queries := sqlc.New(connPool)

	return &DB{connPool, queries}
}
