package postgres

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"social-network/src/pkg/logger"
	"time"
)

type Config struct {
	Host     string
	Port     string
	UserName string
	UserPass string
	DBName   string
	SSLMode  string
}

func MustCreatePostgresDB(ctx context.Context, cfg Config, logger logger.Logger) *pgxpool.Pool {
	var (
		pool   *pgxpool.Pool
		err    error
		config *pgxpool.Config
	)
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.UserName, cfg.UserPass, cfg.Host, cfg.Port, cfg.DBName)

	err = doWithTries(func() error {
		ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx1, url)
		if err != nil {
			return err
		}

		config, err = pgxpool.ParseConfig(url)
		if err != nil {
			return err
		}

		config.ConnConfig.Tracer = NewPostgresLogger(logger)
		pool, err = pgxpool.NewWithConfig(ctx1, config)
		if err != nil {
			return err
		}

		if err = pool.Ping(ctx1); err != nil {
			return err
		}

		return nil
	}, 20, 3*time.Second)

	if err != nil {
		logger.Fatal("Error do with tries postgresql")
	}

	logger.Info("Creating tables for postgres")
	err = createTablesAndIndexes(ctx, pool)
	if err != nil {
		logger.Fatal(err.Error())
	} else {
		logger.Info("All tables for postgres created")
	}

	return pool
}

//go:embed queries/up.sql
var upQuery string

func createTablesAndIndexes(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, upQuery)
	return err
}

func doWithTries(fn func() error, attempts uint8, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return err
}
