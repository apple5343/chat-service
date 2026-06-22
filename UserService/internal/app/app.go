package app

import (
	"context"
	"fmt"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"time"
	"user-service/internal/closer"
	"user-service/internal/config"
	"user-service/internal/server/interceptor"
	authpb "user-service/pkg/api/auth_v1"
	userpb "user-service/pkg/api/user_v1"
	"user-service/pkg/zlog"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		//a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initGRPCServer,
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

func (a *App) initGRPCServer(ctx context.Context) error {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "user-service",
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Info().Msgf("Circuit Breaker: %s, changed from %v, to %v\n", name, from, to)
		},
	})

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.NewCircuitBreakerInterceptor(cb).Unary,
			),
		),
	)

	reflection.Register(a.grpcServer)

	userpb.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserServer(ctx))
	authpb.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthServer(ctx))

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	cfg := a.serviceProvider.LoggerConfig()

	log.Logger = zlog.Default(cfg.IsPretty(), cfg.Version(), cfg.LogLevel())

	return nil
}

func (a *App) runGRPCServer() error {
	fmt.Println("GRPC сервер запущен")

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return fmt.Errorf("listen tcp: %w", err)
	}

	if err = a.grpcServer.Serve(list); err != nil {
		return fmt.Errorf("serve grpc: %w", err)
	}

	return nil
}
