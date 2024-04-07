package main

import (
	"backend-trainee-banner-avito/internal/config"
	"fmt"
	"log/slog"
	"os"
)

const (
	LocalEnv = "local"
	DevEnv   = "dev"
	ProdEnv  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg.Env)
	log := setupLogger(cfg.Env)
	log.Info("Application started....", slog.String("env", cfg.Env))
	log.Debug("Debug messages are active")
	// TODO : LOGGER  : slog
	// TODO : STORAGE : postgresql
	// TODO : ROUTER  : chi, chi-render
	// TODO : SERVER

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case LocalEnv:
		log = slog.New(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case DevEnv:
		log = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case ProdEnv:
		log = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
