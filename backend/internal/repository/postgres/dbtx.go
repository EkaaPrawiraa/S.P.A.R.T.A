package postgres

import (
	"context"
	"database/sql"
)

// DBTX abstracts *sql.DB and *sql.Tx for repository methods.
// It enables UnitOfWork to pass a transaction without changing repository APIs.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
