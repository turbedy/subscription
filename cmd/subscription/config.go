package main

import (
	"time"

	"github.com/turbedy/env/v2"
	"github.com/turbedy/subscription/internal/postgres"
)

type Config struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Postgres        *postgres.Config
}

func Load() (*Config, error) {
	cfg := &Config{
		Host:            "0.0.0.0",
		Port:            "8080",
		ReadTimeout:     4 * time.Minute,
		WriteTimeout:    4 * time.Minute,
		ShutdownTimeout: 8 * time.Minute,
		Postgres:        postgres.NewConfig(),
	}

	err := env.Decode(cfg, env.WithPrefix("SUBSCRIPTION"))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
