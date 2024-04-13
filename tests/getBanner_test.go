package tests

import (
	"backend-trainee-banner-avito/internal/config"
	"backend-trainee-banner-avito/internal/http-server/handlers/banners"
	"backend-trainee-banner-avito/internal/repositories"
	"backend-trainee-banner-avito/internal/storage/postgres"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type Content struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Text string `json:"text"`
}

func TestGetBannerHandler(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug}))
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "testdb",
		},
	}
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	db, err := postgres.NewPG(context.Background(), connString, log, cfg)
	if err != nil {
		log.Error("Failed to initialize test tables")
	}
	defer db.Close()
	errors := CreateTestTablesAndDate(context.Background(), db, log, cfg)
	if errors != nil {
		log.Error("Error with creating db")
	}
	br := repositories.NewBannerRepository(db.Db, log)
	handler := banners.NewGetBannerHandler(log, br)
	req, err := http.NewRequest("GET", "/user_banner?feature_id=1&tag_id=1&use_last_revision=false", nil)
	if err != nil {
		log.Error("Failed to create request")
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code")
	var content Content
	err = json.Unmarshal(rr.Body.Bytes(), &content)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
	expectedContent := Content{
		Text: "Test text",
		Url:  "http://github.com",
		Name: "Test banner",
	}
	assert.Equal(t, expectedContent, content, "Unexpected banner content")
}

func CreateTestTablesAndDate(ctx context.Context, db *postgres.Postgres, log *slog.Logger, cfg *config.Config) error {
	err := postgres.CreateTables(ctx, db.Db, log, cfg)
	if err != nil {
		return err
	}
	jsonContent := Content{
		Text: "Test text",
		Url:  "http://github.com",
		Name: "Test banner",
	}
	var featureID int
	err = db.Db.QueryRow(ctx, `INSERT INTO features (name) VALUES ($1) RETURNING id`, "testFeature").Scan(&featureID)
	if err != nil {
		return err
	}

	var tagID int
	err = db.Db.QueryRow(ctx, `INSERT INTO tags (name) VALUES ($1) RETURNING id`, "testTag").Scan(&tagID)
	if err != nil {
		return err
	}

	var bannerID int
	err = db.Db.QueryRow(ctx, `INSERT INTO banners (feature_id, content, is_active) VALUES ($1, $2, $3) RETURNING id`, featureID, jsonContent, true).Scan(&bannerID)
	if err != nil {
		return err
	}

	_, err = db.Db.Exec(ctx, `INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2)`, bannerID, tagID)
	if err != nil {
		return err
	}
	return nil

}
