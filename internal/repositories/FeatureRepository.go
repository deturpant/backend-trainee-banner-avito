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
	_, err := fr.db.Exec(ctx,
		`INSERT INTO features (name)
			VALUES ($1)
						`, feature.Name)
	if err != nil {
		fr.log.Error("Failed to create Feature", errMsg.Err(err))
	}
	return nil
}
func (fr *FeatureRepository) findFeatureByName(ctx context.Context, name string) (entities.Feature, error) {
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
