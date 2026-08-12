package app

import (
	"api-gateway/internal/closer"
	"api-gateway/internal/config"
	"api-gateway/internal/server"
	"api-gateway/pkg/zlog"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

// App структура приложения, содержащая сервис провайдер и http сервер.
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

// NewApp инициализирует зависимости и создаёт новое приложение.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	return a, nil
}

// Run запускает http сервер.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runHTTPServer()
}

// Инициализация зависимостей.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		//a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return fmt.Errorf("init deps: %w", err)
		}
	}

	return nil
}

// Инициализация конфига.
func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	return nil
}

// Инициализация сервис провайдера.
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

// Инициализация логгера.
func (a *App) initLogger(_ context.Context) error {
	cfg := a.serviceProvider.LoggerConfig()

	log.Logger = zlog.Default(cfg.IsPretty(), cfg.Version(), cfg.LogLevel())

	return nil
}

// Инициализация http сервера.
func (a *App) initHTTPServer(_ context.Context) error {
	srv := server.NewServer(a.serviceProvider.UserHandler(), a.serviceProvider.AuthHandler(), a.serviceProvider.ProjectHandler(), a.serviceProvider.TaskHandler(), a.serviceProvider.ChatHandler())

	httpServer := &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: srv,
	}
	a.httpServer = httpServer

	return nil
}

// Запуск http сервера.
func (a *App) runHTTPServer() error {
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Error starting http server")
		os.Exit(1)
	}
	return nil
}
