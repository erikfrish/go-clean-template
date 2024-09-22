package postgres

import (
	"context"
	"errors"
	"net/url"

	"go-clean-template/config"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPool(cfg config.DB) (*pgxpool.Pool, error) {
	if !cfg.Enabled {
		emptyPool := pgxpool.Pool{}
		return &emptyPool, nil
	}

	query := url.Values{}
	switch cfg.Scheme {
	case "postgres":
		query.Add("dbname", cfg.Database)
		if !cfg.SSLMode {
			query.Add("sslmode", "disable")
		}
	case "sqlserver":
		query.Add("database", cfg.Database)
		query.Add("failoverpartner", cfg.FailoverHost)
	case "clickhouse":
		query.Add("database", cfg.Database)
		query.Add("username", cfg.Username)
		query.Add("password", cfg.Password)
		if !cfg.SSLMode {
			query.Add("sslmode", "disable")
		}
	default:
		return nil, errors.New("unknown db scheme")
	}

	host := cfg.Host
	if cfg.Port != "" {
		host += ":" + cfg.Port
	}

	u := &url.URL{
		Scheme:   cfg.Scheme,
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     host,
		RawQuery: query.Encode(),
	}

	pool, err := pgxpool.New(context.Background(), u.String())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
