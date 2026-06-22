package main

import (
	"context"
	"notification-service/internal/app"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	time.Sleep(time.Duration(4) * time.Second)

	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to init app")
		os.Exit(1)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to run app")
		os.Exit(1)
	}
}
