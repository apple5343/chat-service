package config

import (
	"errors"
	"net"
	"os"
)

var (
	errHTTPHostNotFound = errors.New("http host not found")
	errHTTPPortNotFound = errors.New("http port not found")
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

// NewHTTPConfig инициализирует http конфиг.
func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errHTTPHostNotFound
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errHTTPPortNotFound
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
