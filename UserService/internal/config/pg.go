package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	errPGUserNotFound = errors.New("pg user not found")
	errPGPassNotFound = errors.New("pg pass not found")
	errPGHostNotFound = errors.New("pg host not found")
	errPGPortNotFound = errors.New("pg port not found")
	errPGDBNotFound   = errors.New("pg db not found")
)

const (
	pgUserEnvName = "PG_USER"
	pgPassEnvName = "PG_PASS"
	pgHostEnvName = "PG_HOST"
	pgPortEnvName = "PG_PORT"
	pgDBEnvName   = "PG_DB"
)

type PGConfig interface {
	PGConn() string
}

type pgConfig struct {
	pgUser string
	pgPass string
	pgHost string
	pgPort string
	pgDB   string
}

func NewPGConfig() (PGConfig, error) {
	pgUser := os.Getenv(pgUserEnvName)
	if len(pgUser) == 0 {
		return nil, errPGUserNotFound
	}
	pgPass := os.Getenv(pgPassEnvName)
	if len(pgPass) == 0 {
		return nil, errPGPassNotFound
	}
	pgHost := os.Getenv(pgHostEnvName)
	if len(pgHost) == 0 {
		return nil, errPGHostNotFound
	}
	pgPort := os.Getenv(pgPortEnvName)
	if len(pgPort) == 0 {
		return nil, errPGPortNotFound
	}
	pgDB := os.Getenv(pgDBEnvName)
	if len(pgDB) == 0 {
		return nil, errPGDBNotFound
	}

	return &pgConfig{
		pgUser: pgUser,
		pgPass: pgPass,
		pgHost: pgHost,
		pgPort: pgPort,
		pgDB:   pgDB,
	}, nil
}

func (cfg *pgConfig) PGConn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.pgUser, cfg.pgPass, cfg.pgHost, cfg.pgPort, cfg.pgDB)
}
