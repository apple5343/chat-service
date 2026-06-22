package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

// Load подгружает .env файл.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	return nil
}
