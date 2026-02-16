package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres wraps a PostgreSQL connection pool as the primary database.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres creates a new PostgreSQL connection pool.
// Caller must call Close when done.
func NewPostgres(ctx context.Context, connStr string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

// NewPostgresOrNil is like NewPostgres but logs and returns nil on error so the app can run without DB.
func NewPostgresOrNil(ctx context.Context, connStr string) *Postgres {
	p, err := NewPostgres(ctx, connStr)
	if err != nil {
		log.Printf("warning: could not create postgres pool: %v", err)
		return nil
	}
	if err := p.Ping(ctx); err != nil {
		log.Printf("warning: postgres ping failed: %v", err)
	}
	return p
}

// Ping checks connectivity to PostgreSQL.
func (p *Postgres) Ping(ctx context.Context) error {
	if p == nil || p.pool == nil {
		return nil
	}
	return p.pool.Ping(ctx)
}

// Close closes the connection pool.
func (p *Postgres) Close() {
	if p != nil && p.pool != nil {
		p.pool.Close()
	}
}

// Pool returns the underlying pgx pool for running queries. May be nil.
func (p *Postgres) Pool() *pgxpool.Pool {
	if p == nil {
		return nil
	}
	return p.pool
}
