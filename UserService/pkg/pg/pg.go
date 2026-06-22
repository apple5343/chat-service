package pg

import (
	"github.com/rs/zerolog/log"
	"os"

	"github.com/jmoiron/sqlx"
)

// Init инициализирует Postgres db.
func Init(pgURL string) *sqlx.DB {
	db, err := sqlx.Connect("pgx", pgURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to db")
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		log.Error().Err(err).Msg("failed to ping db")
		os.Exit(1)
	}

	return db
}
