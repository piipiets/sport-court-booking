package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/piipiets/sport-court-booking/middlewares"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/model/dto/response"
	"github.com/piipiets/sport-court-booking/repository"
)

// Service interface
type Service interface {
	Login(ctx context.Context, req request.LoginRequest) (response.LoginResponse, error)
}

type userService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) Service {
	return &userService{
		repo: repo,
	}
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already registered")
)

func (s *userService) Login(ctx context.Context, req request.LoginRequest) (response.LoginResponse, error) {
	// 1. Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Jika user not found, return error yang sama dengan password salah (security)
		if errors.Is(err, ErrUserNotFound) {
			return response.LoginResponse{}, errors.New("invalid email or password")
		}
		return response.LoginResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Verify password dengan bcrypt
	if !req.VerifyPassword(user.Password, req.Password) {
		return response.LoginResponse{}, errors.New("invalid email or password")
	}

	// 3. Generate JWT token
	jwtToken, err := middlewares.GenerateJwtToken()
	if err != nil {
		return response.LoginResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	// 4. Return success response
	return response.LoginResponse{
		Token:     jwtToken,
		ExpiredAt: time.Now().Add(1 * time.Minute),
	}, nil
}
