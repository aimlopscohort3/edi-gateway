package config

import "os"

type Config struct {
    DatabaseDSN     string
    KafkaBroker     string
    TemporalHostPort string
}

func Load() Config {
    return Config{
        DatabaseDSN:     getEnv("DATABASE_DSN", "host=postgres user=postgres password=postgres dbname=edi_gateway port=5432 sslmode=disable"),
        KafkaBroker:     getEnv("KAFKA_BROKER", "broker:9092"),
        TemporalHostPort: getEnv("TEMPORAL_HOST_PORT", "localhost:7233"),
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
