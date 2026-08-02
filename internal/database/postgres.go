package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/config"
)

func NewPostgresPool(
	lifecycle fx.Lifecycle,
	applicationConfig config.Config,
) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(
		applicationConfig.DatabaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		poolConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			if err := pool.Ping(context); err != nil {
				pool.Close()
				return fmt.Errorf("ping PostgreSQL: %w", err)
			}

			log.Println("Connected to PostgreSQL")
			return nil
		},
		OnStop: func(context.Context) error {
			log.Println("Closing PostgreSQL pool")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
