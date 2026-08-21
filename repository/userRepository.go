package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/piipiets/sport-court-booking/model/entity"
)

// Error definitions
var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already registered")
)

// Repository interface
type Repository interface {
	SignUp(ctx context.Context, user entity.User) error
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
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
		if isDuplicateKeyError(err) {
			return ErrEmailExists
		}
		return fmt.Errorf("failed to create user: %w", err)
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
			return entity.User{}, ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// CheckEmailExists checks if email already registered
func (r *userRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// isDuplicateKeyError checks if error is duplicate key violation
func isDuplicateKeyError(err error) bool {
	errMsg := err.Error()
	return containsAny(errMsg, []string{
		"duplicate key",
		"Duplicate entry",
		"UNIQUE constraint failed",
		"unique constraint",
	})
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) && s[:len(substr)] == substr {
			return true
		}
	}
	return false
}
