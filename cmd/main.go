package main

import (
	"backend-trainee-banner-avito/internal/config"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"backend-trainee-banner-avito/internal/storage/postgres"
	"context"
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
	log.Debug("Debug messages are active")

	//CONNECT TO DB
	pg, err := connectToPostgres(cfg, log)
	if err != nil {
		log.Error("Error creating PostgreSQL instance", errMsg.Err(err))
		os.Exit(1)
	}
	defer pg.Close()

	if err := pg.Ping(context.Background()); err != nil {
		log.Error("Error pinging PostgreSQL:", errMsg.Err(err))
		os.Exit(1)
	} else {
		log.Info("Connected to PostgreSQL successfully!")
	}
	log.Info("Application started....", slog.String("env", cfg.Env))

	/*ur := repositories.NewRepository(pg.Db, log)
	errorr := ur.CreateUser(context.Background(), &entities.User{Role: "user", Password: "123123", Username: "pop"})
	if err != nil {
		log.Error("error adding user", errMsg.Err(errorr))
	}*/

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

func connectToPostgres(cfg *config.Config, log *slog.Logger) (*postgres.Postgres, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	pg, err := postgres.NewPG(context.Background(), connString, log)
	return pg, err
}
