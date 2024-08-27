package auth

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/qRe0/auth-api/configs"
	"github.com/qRe0/auth-api/internal/models"

	_ "github.com/lib/pq"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func Init(cfg configs.DBConfig) (*sqlx.DB, error) {
	const connectionStringTemplate = "postgres://%s:%s@%s:%s/%s?sslmode=disable"

	connectionString := fmt.Sprintf(connectionStringTemplate, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres db")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping postgres db")
	}

	return db, nil
}

func (a *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	const createUserQuery = `INSERT INTO users (name, phone, email, password) VALUES ($1, $2, $3, $4)`

	_, err := a.db.ExecContext(ctx, createUserQuery, user.Name, user.Phone, user.Email, user.Password)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (a *AuthRepository) GetUserData(ctx context.Context, phone string) (models.User, error) {
	const getUserIDQuery = `SELECT user_id FROM Users WHERE phone = $1`
	var data models.User
	err := a.db.QueryRowxContext(ctx, getUserIDQuery, phone).Scan(&data.ID)
	if err != nil {
		return models.User{}, err
	}

	const getUserPassQuery = `SELECT password FROM Users WHERE phone = $1`
	err = a.db.QueryRowxContext(ctx, getUserPassQuery, phone).Scan(&data.Password)
	if err != nil {
		return models.User{}, err
	}

	return data, nil
}
