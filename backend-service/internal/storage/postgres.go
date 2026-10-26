package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"semen_project/internal/config"
	"semen_project/internal/errs"
)

// ConnectToPg создает подключение к PostgreSQL
func ConnectToPg(ctx context.Context, cfg config.DbConfig, appName string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.PgConnDsn())
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errs.ErrWrongConfig, err.Error())
	}

	pgxConfig.MaxConns = int32(cfg.PgMaxConn)
	pgxConfig.ConnConfig.RuntimeParams["application_name"] = appName

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errs.ErrEstablishConnect, err.Error())
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%w: ping failed: %s", errs.ErrEstablishConnect, err.Error())
	}

	return pool, nil
}
