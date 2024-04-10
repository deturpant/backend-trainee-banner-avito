package repositories

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type BannerTagRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewBannerTagRepository(db *pgxpool.Pool, log *slog.Logger) *BannerTagRepository {
	return &BannerTagRepository{db, log}
}

func (btr *BannerTagRepository) CreateBannerTag(ctx context.Context, bannerTag *entities.BannerTag) error {
	_, err := btr.db.Exec(ctx,
		`INSERT INTO banner_tags (banner_id, tag_id)
			VALUES ($1, $2)`, bannerTag.BannerID, bannerTag.TagID)
	if err != nil {
		btr.log.Error("Failed to create BannerTag", errMsg.Err(err))
		return err
	}
	return nil
}

func (btr *BannerTagRepository) FindBannerTagsByBannerID(ctx context.Context, bannerID int) ([]entities.BannerTag, error) {
	rows, err := btr.db.Query(ctx, `SELECT * FROM banner_tags WHERE banner_id = $1`, bannerID)
	if err != nil {
		btr.log.Error("Failed to find BannerTags by Banner ID", errMsg.Err(err))
		return nil, err
	}
	defer rows.Close()

	var bannerTags []entities.BannerTag
	for rows.Next() {
		var bannerTag entities.BannerTag
		if err := rows.Scan(&bannerTag.BannerID, &bannerTag.TagID); err != nil {
			btr.log.Error("Failed to scan BannerTag row", errMsg.Err(err))
			return nil, err
		}
		bannerTags = append(bannerTags, bannerTag)
	}
	if err := rows.Err(); err != nil {
		btr.log.Error("Error occurred while iterating BannerTag rows", errMsg.Err(err))
		return nil, err
	}
	return bannerTags, nil
}
