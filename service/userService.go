package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/piipiets/sport-court-booking/middlewares"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/model/dto/response"
	"github.com/piipiets/sport-court-booking/model/entity"
	"github.com/piipiets/sport-court-booking/repository"
)

// Service interface
type UserService interface {
	Login(ctx context.Context, req request.LoginRequest) (response.LoginResponse, error)
	SignUp(ctx context.Context, req request.SignUpRequest) error
}

type userService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) UserService {
	return &userService{
		repo: repo,
	}
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already registered")
)

func (s *userService) Login(ctx context.Context, req request.LoginRequest) (response.LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return response.LoginResponse{}, errors.New("invalid email or password")
		}
		return response.LoginResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password with bcrypt
	if !req.VerifyPassword(user.Password, req.Password) {
		return response.LoginResponse{}, errors.New("invalid email or password")
	}

	// Generate JWT token
	jwtToken, err := middlewares.GenerateJwtToken(user)
	if err != nil {
		return response.LoginResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return response.LoginResponse{
		Token:     jwtToken,
		ExpiredAt: time.Now().Add(10 * time.Minute),
	}, nil
}

func (s *userService) SignUp(ctx context.Context, req request.SignUpRequest) (err error) {

	var user entity.User

	user, err = req.ConvertToModelForSignUp()
	if err != nil {
		return err
	}

	err = s.repo.SignUp(ctx, user)
	if err != nil {
		fmt.Println("Error sign up : ", err)
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return errors.New("user already exists")
			}
		}

		return fmt.Errorf("failed to create user")
	}

	return nil
}
