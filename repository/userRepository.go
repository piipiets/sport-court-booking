package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/model/entity"
)

// Repository interface
type Repository interface {
	SignUp(ctx context.Context, user entity.User) error
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) Repository {
	return &userRepository{
		db: db,
	}
}

// SignUp creates a new user
func (r *userRepository) SignUp(ctx context.Context, user entity.User) error {
	var userID int64

	queryInsertUser := `INSERT INTO users (email, name, password, role, created_at) 
						VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, queryInsertUser,
		user.Email,
		user.Name,
		user.Password,
		user.Role,
		time.Now(),
	).Scan(&userID)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail retrieves user by email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	queryGetUserByEmail := `SELECT id, email, name, password, role, created_at FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, queryGetUserByEmail, email).
		Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Role, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, constant.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}
