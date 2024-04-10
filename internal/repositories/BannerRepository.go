package repositories

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type BannerRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewBannerRepository(db *pgxpool.Pool, log *slog.Logger) *BannerRepository {
	return &BannerRepository{db, log}
}

func (br *BannerRepository) CreateBanner(ctx context.Context, banner *entities.Banner) error {
	_, err := br.db.Exec(ctx,
		`INSERT INTO banners (feature_id, content, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)`,
		banner.FeatureID, banner.Content, banner.IsActive, banner.CreatedAt, banner.UpdatedAt)
	if err != nil {
		br.log.Error("Failed to create banner", errMsg.Err(err))
		return err
	}
	return nil
}

func (br *BannerRepository) FindBannerById(ctx context.Context, id int) (entities.Banner, error) {
	query, err := br.db.Query(ctx,
		`SELECT * FROM banners WHERE id = $1`, id)
	if err != nil {
		br.log.Error("Banner not found", errMsg.Err(err))
		return entities.Banner{}, err
	}
	defer query.Close()
	rowArray := entities.Banner{}
	for query.Next() {
		err := query.Scan(&rowArray.ID, &rowArray.FeatureID, &rowArray.Content, &rowArray.IsActive, &rowArray.CreatedAt, &rowArray.UpdatedAt)
		if err != nil {
			br.log.Error("Banner not found", errMsg.Err(err))
			return entities.Banner{}, err
		}
	}
	return rowArray, nil
}

func (br *BannerRepository) FindBannersByFeatureID(ctx context.Context, feature_id int) ([]entities.Banner, error) {
	query, err := br.db.Query(ctx, `SELECT * FROM banners WHERE feature_id = $1`, feature_id)
	if err != nil {
		br.log.Error("Error querying banners", errMsg.Err(err))
		return nil, err
	}
	defer query.Close()

	var resultSlice []entities.Banner
	for query.Next() {
		var rowArray entities.Banner
		err := query.Scan(&rowArray.ID, &rowArray.FeatureID, &rowArray.Content, &rowArray.IsActive, &rowArray.CreatedAt, &rowArray.UpdatedAt)
		if err != nil {
			br.log.Error("Error scanning banners", errMsg.Err(err))
			return nil, err
		}
		resultSlice = append(resultSlice, rowArray)
	}

	if len(resultSlice) == 0 {
		br.log.Info("No banners found for feature ID:", feature_id)
		return []entities.Banner{}, nil
	}

	return resultSlice, nil
}

func (br *BannerRepository) findBannersByTagID(ctx context.Context, tagId int) ([]entities.Banner, error) {
	query, err := br.db.Query(ctx, `SELECT * FROM banner_tags WHERE tag_id = $1`, tagId)
	if err != nil {
		br.log.Error("Error querying banners", errMsg.Err(err))
		return nil, err
	}
	var resultSlice []entities.Banner // Инициализация пустого слайса
	defer query.Close()
	for query.Next() {
		var rowArray entities.Banner
		err := query.Scan(&rowArray.ID, &rowArray.FeatureID, &rowArray.Content, &rowArray.IsActive, &rowArray.CreatedAt, &rowArray.UpdatedAt)
		if err != nil {
			br.log.Error("Error scanning banners", errMsg.Err(err))
			return nil, err
		}
		resultSlice = append(resultSlice, rowArray)
	}
	if len(resultSlice) == 0 {
		br.log.Info("No banners found for tag ID:", tagId)
		return []entities.Banner{}, nil
	}
	return resultSlice, nil
}
