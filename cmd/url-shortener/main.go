package main

import (
	"log/slog"
	"os"

	"github.com/marchuknikolay/url-shortener/internal/config"
	"github.com/marchuknikolay/url-shortener/internal/config/lib/logger/sl"
	"github.com/marchuknikolay/url-shortener/internal/config/storage/sqlite"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := initLogger(cfg.Env)

	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	storage, err := sqlite.New(cfg.Storage)
	if err != nil {
		log.Error("Error opening a storage, ", sl.Err(err))
		os.Exit(1)
	}

	err = storage.DeleteUrl("google")
	if err != nil {
		log.Error("Error deleting a url, ", sl.Err(err))
		os.Exit(1)
	}
	log.Info("Deleted a url")

	url, err := storage.GetUrl("google")
	if err != nil {
		log.Error("Error getting a url, ", sl.Err(err))
		os.Exit(1)
	}
	log.Info("", slog.String("url", url))

	// init router

	// run server
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
