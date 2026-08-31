package main

import (
	"github.com/turbedy/env/v2"
	"github.com/turbedy/subscription/internal/postgres"
)

type Config struct {
	Postgres *postgres.Config
}

func Load() (*Config, error) {
	cfg := &Config{
		Postgres: postgres.NewConfig(),
	}

	err := env.Decode(cfg, env.WithPrefix("SUBSCRIPTION"))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
