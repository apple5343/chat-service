package config_test

import (
	"github.com/stretchr/testify/require"
	"notification-service/internal/config"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	path := "../../.env"

	err := config.Load(path)
	require.NoError(t, err)
}

func TestLoadConfigError(t *testing.T) {
	path := ""

	err := config.Load(path)
	require.Error(t, err)
}
