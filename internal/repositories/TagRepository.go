package repositories

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type TagRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewTagRepository(db *pgxpool.Pool, log *slog.Logger) *TagRepository {
	return &TagRepository{db, log}
}

func (tr *TagRepository) createTag(ctx context.Context, tag *entities.Tag) error {
	_, err := tr.db.Exec(ctx,
		`INSERT INTO tags (name)
		VALUES ($1)`, tag.Name)
	if err != nil {
		tr.log.Error("Failed to create tag", errMsg.Err(err))
		return err
	}
	return nil
}
func (tr *TagRepository) findTagByName(ctx context.Context, name string) (entities.Tag, error) {
	query, err := tr.db.Query(ctx,
		`SELECT * FROM tags WHERE name = $1`, name)
	if err != nil {
		tr.log.Error("Tag not found", errMsg.Err(err))
		return entities.Tag{}, err
	}
	row := entities.Tag{}
	if !query.Next() {
		tr.log.Error("Tag not found")
		return entities.Tag{}, err
	} else {
		err := query.Scan(&row.ID, &row.Name)
		if err != nil {
			tr.log.Error("Tag not found", errMsg.Err(err))
			return entities.Tag{}, fmt.Errorf("Tag not found")
		}
	}
	return row, nil
}
