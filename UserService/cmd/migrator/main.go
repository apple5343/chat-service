package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	migrationDirEnvName = "MIGRATION_DIR"
	pgConnEnvName       = "PG_CONN"
)

func main() {
	var migrationDir, connection string

	migrationDir, _ = os.LookupEnv(migrationDirEnvName)
	connection, _ = os.LookupEnv(pgConnEnvName)

	if migrationDir == "" || connection == "" {
		log.Println("MIGRATION_DIR and PG_CONN env vars must be set")
		return
	}

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println(err)
		return
	}

	fsrc, err := (&file.File{}).Open(fmt.Sprintf("file://%s", migrationDir))
	if err != nil {
		log.Println(err)
		return
	}

	m, err := migrate.NewWithInstance("file", fsrc, "postgres", driver)
	if err != nil {
		log.Println(err)
		return
	}
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no changes")
			return
		}
		log.Println(err)
		return
	}

	log.Println("success")
}
