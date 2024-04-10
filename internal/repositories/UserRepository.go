package repositories

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type UserRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *UserRepository {
	return &UserRepository{db, log}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	err := ur.db.QueryRow(ctx,
		`INSERT INTO users (username, password, role)
		VALUES ($1,$2,$3)
		RETURNING id`, user.Username, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		ur.log.Error("Failed to create user", errMsg.Err(err))
		return err
	}
	return nil
}
func (ur *UserRepository) FindUserByUsername(ctx context.Context, username string) (entities.User, error) {
	query, err := ur.db.Query(ctx, `SELECT * FROM users WHERE username = $1`, username)
	if err != nil {
		ur.log.Error("Error querying users", errMsg.Err(err))
		return entities.User{}, err
	}
	row := entities.User{}
	defer query.Close()
	if !query.Next() {
		ur.log.Error("User not found")
		return entities.User{}, fmt.Errorf("User not found")
	} else {
		err := query.Scan(&row.ID, &row.Username, &row.Password, &row.Role)
		if err != nil {
			ur.log.Error("Error scanning users", errMsg.Err(err))
			return entities.User{}, err
		}
	}
	return row, nil
}
func (ur *UserRepository) FindUserById(ctx context.Context, id int) (entities.User, error) {
	query, err := ur.db.Query(ctx,
		`SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		ur.log.Error("Error querying users", errMsg.Err(err))
		return entities.User{}, err
	}
	defer query.Close()
	rowArray := entities.User{}
	if !query.Next() {
		ur.log.Error("User not found")
		return entities.User{}, fmt.Errorf("User not found")

	} else {
		err := query.Scan(&rowArray.ID, &rowArray.Username, &rowArray.Password, &rowArray.Role)
		if err != nil {
			ur.log.Error("Error scanning users", errMsg.Err(err))
			return entities.User{}, err
		}
	}
	return rowArray, nil
}
