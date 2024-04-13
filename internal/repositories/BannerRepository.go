package repositories

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/http-server/handlers/banners"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"strconv"
)

type BannerRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewBannerRepository(db *pgxpool.Pool, log *slog.Logger) *BannerRepository {
	return &BannerRepository{db, log}
}

func (br *BannerRepository) CreateBanner(ctx context.Context, banner *entities.Banner) error {
	err := br.db.QueryRow(ctx,
		`INSERT INTO banners (feature_id, content, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		banner.FeatureID, banner.Content, banner.IsActive, banner.CreatedAt, banner.UpdatedAt).Scan(&banner.ID)
	if err != nil {
		br.log.Error("Failed to create banner", errMsg.Err(err))
		return err
	}
	return nil
}

func (br *BannerRepository) FindBannerById(ctx context.Context, id int) (entities.Banner, error) {
	row := br.db.QueryRow(ctx,
		`SELECT * FROM banners WHERE id = $1`, id)

	var banner entities.Banner
	err := row.Scan(&banner.ID, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если баннер не найден, возвращаем ошибку
			return entities.Banner{}, err
		}
		br.log.Error("Failed to scan banner row", errMsg.Err(err))
		return entities.Banner{}, err
	}

	return banner, nil
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

func (br *BannerRepository) FindBannerByFeatureTag(ctx context.Context, featureID, tagID int) (*entities.Banner, error) {
	query := `SELECT b.id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
			  FROM banners b
			  INNER JOIN banner_tags bt ON b.id = bt.banner_id
			  WHERE b.feature_id = $1 AND bt.tag_id = $2 AND b.is_active = true`
	row := br.db.QueryRow(ctx, query, featureID, tagID)
	var banner entities.Banner
	err := row.Scan(&banner.ID, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
	if err != nil {
		br.log.Error("Error with database", errMsg.Err(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		br.log.Error("Failed to find banner", errMsg.Err(err))
		return nil, err
	}
	return &banner, nil
}

func (br *BannerRepository) DeleteBannerByID(ctx context.Context, id int) error {
	_, err := br.db.Exec(ctx, `DELETE FROM banners WHERE id = $1`, id)
	if err != nil {
		br.log.Error("Failed to delete banner by ID", errMsg.Err(err))
		return err
	}
	return nil
}

func (br *BannerRepository) FindBannersByParameters(ctx context.Context, params banners.RequestGetBanners) ([]entities.Banner, error) {
	query := "SELECT b.id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at, array_agg(bt.tag_id) AS tag_ids FROM banners b LEFT JOIN banner_tags bt ON b.id = bt.banner_id WHERE 1=1"
	args := []interface{}{}

	if params.FeatureID != nil {
		query += " AND b.feature_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, *params.FeatureID)
	}

	if params.TagID != nil {
		query += " AND b.id IN (SELECT banner_id FROM banner_tags WHERE tag_id = $" + strconv.Itoa(len(args)+1) + ")"
		args = append(args, *params.TagID)
	}

	query += " GROUP BY b.id"

	if params.Limit != nil {
		query += " LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, *params.Limit)
	}

	if params.Offset != nil {
		query += " OFFSET $" + strconv.Itoa(len(args)+1)
		args = append(args, *params.Offset)
	}

	rows, err := br.db.Query(ctx, query, args...)
	if err != nil {
		br.log.Error("Failed to query banners", errMsg.Err(err))
		return nil, err
	}
	defer rows.Close()

	var banners []entities.Banner
	for rows.Next() {
		var banner entities.Banner
		var tagIDs []int
		if err := rows.Scan(&banner.ID, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt, &tagIDs); err != nil {
			br.log.Error("Failed to scan banner row", errMsg.Err(err))
			return nil, err
		}
		banner.TagIDs = tagIDs
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		br.log.Error("Error occurred while iterating banner rows", errMsg.Err(err))
		return nil, err
	}

	return banners, nil
}

func (br *BannerRepository) UpdateBanner(ctx context.Context, banner *entities.Banner) error {
	tx, err := br.db.Begin(ctx)
	if err != nil {
		br.log.Error("Failed to begin transaction", errMsg.Err(err))
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM banner_tags WHERE banner_id = $1`, banner.ID)
	if err != nil {
		br.log.Error("Failed to delete old tags for banner", errMsg.Err(err))
		return err
	}

	for _, tagID := range banner.TagIDs {
		_, err = tx.Exec(ctx, `INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2)`, banner.ID, tagID)
		if err != nil {
			br.log.Error("Failed to insert tag for banner", errMsg.Err(err))
			return err
		}
	}

	_, err = tx.Exec(ctx,
		`UPDATE banners SET feature_id = $1, content = $2, is_active = $3, updated_at = $4 WHERE id = $5`,
		banner.FeatureID, banner.Content, banner.IsActive, banner.UpdatedAt, banner.ID)
	if err != nil {
		br.log.Error("Failed to update banner", errMsg.Err(err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		br.log.Error("Failed to commit transaction", errMsg.Err(err))
		return err
	}

	return nil
}
