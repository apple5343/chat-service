package main

import (
	"context"
	"os"
	"project-service/internal/app"
	"time"

	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	time.Sleep(time.Duration(4) * time.Second)

	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to init app")
		os.Exit(1)
	}

	err = a.Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to run app")
		os.Exit(1)
	}
}
