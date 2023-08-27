package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lekht/notepad/config"
	"github.com/pkg/errors"
)

type PostgreDB struct {
	Pool *pgxpool.Pool
}

func New(cfg *config.Postgres) (*PostgreDB, error) {
	var connstr string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgPort, cfg.PgDB)

	pool, err := pgxpool.New(context.Background(), connstr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make connection pool")
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping postgres")
	}

	return &PostgreDB{pool}, err
}
