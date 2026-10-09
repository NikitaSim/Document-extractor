package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr      string
	HTTPTimeout   time.Duration
	MaxFileSize   int64
	OllamaURL     string
	OllamaModel   string
	OllamaTimeout time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:      env("HTTP_ADDR", ":8080"),
		HTTPTimeout:   duration("HTTP_TIMEOUT", 5*time.Minute),
		MaxFileSize:   int64Value("MAX_FILE_SIZE", 20*1024*1024),
		OllamaURL:     env("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel:   env("OLLAMA_MODEL", "qwen3:8b"),
		OllamaTimeout: duration("OLLAMA_TIMEOUT", 5*time.Minute),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func duration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}

func int64Value(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if n, err := strconv.ParseInt(value, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}
