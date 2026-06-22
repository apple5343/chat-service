package config

import (
	"errors"
	"net"
	"os"
)

var (
	errGRPCHostNotFound = errors.New("grpc host not found")
	errGRPCPortNotFound = errors.New("grpc port not found")
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errGRPCHostNotFound
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errGRPCPortNotFound
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
