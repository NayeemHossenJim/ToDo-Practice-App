package database

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/db"
)

func NewQuerier(pool *pgxpool.Pool) db.Querier {
	return db.New(pool)
}
