package app

import (
	"context"
	"fmt"
	"notification-service/internal/closer"
	"notification-service/internal/config"
	"notification-service/pkg/zlog"

	"github.com/rs/zerolog/log"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runKafkaConsumer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		//a.initConfig,
		a.initServiceProvider,
		a.initLogger,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return fmt.Errorf("init deps: %w", err)
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	cfg := a.serviceProvider.LoggerConfig()

	log.Logger = zlog.Default(cfg.IsPretty(), cfg.Version(), cfg.LogLevel())

	return nil
}

func (a *App) runKafkaConsumer(ctx context.Context) error {
	for {
		if err := a.serviceProvider.ConsumerGroup().Consume(ctx, a.serviceProvider.KafkaConfig().Topics(), a.serviceProvider.Consumer()); err != nil {
			log.Error().Err(err).Msgf("Error in consuming: %v", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		a.serviceProvider.Consumer().Prepare()
	}
}
