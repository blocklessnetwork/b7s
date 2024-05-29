package telemetry

import (
	"time"
)

type Option func(*Config)

type Config struct {
	ID                string
	ExporterMethod    ExporterMethod
	BatchTraceTimeout time.Duration
}

var defaultConfig = Config{
	ExporterMethod:    ExporterGRPC,
	BatchTraceTimeout: 1 * time.Second,
}

func WithExporterMethod(m ExporterMethod) Option {
	return func(cfg *Config) {
		cfg.ExporterMethod = m
	}
}

func WithBatchTraceTimeout(t time.Duration) Option {
	return func(cfg *Config) {
		cfg.BatchTraceTimeout = t
	}
}

func WithID(id string) Option {
	return func(cfg *Config) {
		cfg.ID = id
	}
}
