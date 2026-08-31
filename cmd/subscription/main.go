package main

import (
	"log/slog"
	"os"

	"github.com/turbedy/subscription/internal/postgres"
)

var Version = "dev"

func main() {
	log := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	))
	log.Info("starting subscription", "version", Version)

	cfg, err := Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := postgres.Open(cfg.Postgres)
	if err != nil {
		log.Error("failed to open postgres", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("failed to close postgres", "error", err)
		}
	}()
}
