package main

import (
	"backend-trainee-banner-avito/internal/config"
	"backend-trainee-banner-avito/internal/http-server/handlers/features"
	"backend-trainee-banner-avito/internal/http-server/handlers/tags"
	"backend-trainee-banner-avito/internal/http-server/handlers/users"
	"backend-trainee-banner-avito/internal/lib/auth"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"backend-trainee-banner-avito/internal/repositories"
	"backend-trainee-banner-avito/internal/storage/postgres"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
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

	//router init
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	fr := repositories.NewFeatureRepository(pg.Db, log)
	router.Post("/features", features.New(log, fr))
	tr := repositories.NewTagRepository(pg.Db, log)
	ur := repositories.NewUserRepository(pg.Db, log)
	router.Post("/tags", tags.New(log, tr))
	router.Post("/users", users.New(log, ur))
	jwt := auth.NewJWTManager("sdmgslgnjfkd", log)
	router.Post("/login", users.LoginFunc(log, ur, jwt))

	log.Info("Starting server at", slog.String("addr", cfg.Server.Addr))
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
	}

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
