package repositories

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type FeatureRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewFeatureRepository(db *pgxpool.Pool, log *slog.Logger) *FeatureRepository {
	return &FeatureRepository{db, log}
}

func (fr *FeatureRepository) CreateFeature(ctx context.Context, feature *entities.Feature) error {
	err := fr.db.QueryRow(ctx,
		`INSERT INTO features (name)
		VALUES ($1)
		RETURNING id`, feature.Name).Scan(&feature.ID)
	if err != nil {
		fr.log.Error("Failed to create feature", errMsg.Err(err))
		return err
	}
	return nil
}
func (fr *FeatureRepository) FindFeatureById(ctx context.Context, id int) (entities.Feature, error) {
	var row entities.Feature
	err := fr.db.QueryRow(ctx, `SELECT id, name FROM features WHERE id = $1`, id).Scan(&row.ID, &row.Name)
	if err != nil {
		fr.log.Error("Failed to find Feature by ID", errMsg.Err(err))
		return entities.Feature{}, err
	}
	return row, nil
}

func (fr *FeatureRepository) FindFeatureByName(ctx context.Context, name string) (entities.Feature, error) {
	query, err := fr.db.Query(ctx, `SELECT * FROM features WHERE name = $1`, name)
	if err != nil {
		fr.log.Error("Feature not found", errMsg.Err(err))
		return entities.Feature{}, err
	}
	row := entities.Feature{}
	if !query.Next() {
		fr.log.Error("Feature not found")
		return entities.Feature{}, fmt.Errorf("Feature not found")
	} else {
		err := query.Scan(&row.ID, &row.Name)
		if err != nil {
			fr.log.Error("Feature not found", errMsg.Err(err))
		}
	}
	return row, nil

}
