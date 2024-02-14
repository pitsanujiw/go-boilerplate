package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pitsanujiw/go-boilerplate/config"
)

func AcquireDBPool(ctx context.Context, cfg config.Database) (*pgxpool.Pool, error) {
	var (
		dbURL = fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s",
			cfg.UserName,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)

		options = make([]string, 0)
	)

	if cfg.Timeout > 0 {
		options = append(options, fmt.Sprintf("connect_timeout=%d", cfg.Timeout))
	}

	if cfg.Schema != "" {
		options = append(options, fmt.Sprintf("schema=%s", cfg.Schema))
	}

	if cfg.SSL {
		options = append(options, "sslmode=prefer")
		if cfg.Certificate != "" {
			// Path of the server certificate. Certificate paths are resolved relative to the
			options = append(options, fmt.Sprintf("sslcert=%s", cfg.Certificate))
		}
	}

	cfgPool, err := pgxpool.ParseConfig(fmt.Sprintf("%s?%s", dbURL, strings.Join(options, "&")))
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfgPool)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
