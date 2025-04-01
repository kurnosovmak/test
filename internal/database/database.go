package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kurnosovmak/test/internal/config"
	"github.com/kurnosovmak/test/pkg/logger"
)

const duplicatedKey = "23505"

type Database struct {
	pool *pgxpool.Pool
}

func IsDuplicatedKeyError(err error, field string) bool {
	var perr *pgconn.PgError
	if !errors.As(err, &perr) {
		return false

	}
	logger.Info(perr.ConstraintName)
	return perr.Code == duplicatedKey && perr.ConstraintName == fmt.Sprintf("%s_key", field)
}

func New(cfg *config.Database) (*Database, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return &Database{pool: pool}, nil
}

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *Database) GetPool() *pgxpool.Pool {
	return db.pool
}
