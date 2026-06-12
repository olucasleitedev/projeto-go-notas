package config

import (
	"os"
	"strings"
)

func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func KafkaBrokers() []string {
	raw := EnvOr("KAFKA_BROKERS", "localhost:19092")
	parts := strings.Split(raw, ",")
	brokers := make([]string, 0, len(parts))
	for _, p := range parts {
		if b := strings.TrimSpace(p); b != "" {
			brokers = append(brokers, b)
		}
	}
	return brokers
}

func KafkaEnabled() bool {
	return EnvOr("KAFKA_ENABLED", "true") != "false"
}
