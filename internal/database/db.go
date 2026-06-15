package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewDB(path string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), path)
}
