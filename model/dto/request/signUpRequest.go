package request

import (
	"errors"
	"time"

	"github.com/piipiets/sport-court-booking/helpers/common"
	entity "github.com/piipiets/sport-court-booking/model/entity"
)

type SignUpRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	ReTypePassword string `json:"re_type_password" binding:"required,eqfield=Password"`
	Name           string `json:"name" binding:"required"`
}

func (s *SignUpRequest) ConvertToModelForSignUp() (user entity.User, err error) {
	hashedPassword, err := common.HashPassword(s.Password)
	if err != nil {
		err = errors.New("hashing password failed")
		return
	}

	return entity.User{
		Email:     s.Email,
		Password:  hashedPassword,
		Name:      s.Name,
		Role:      "user",
		CreatedAt: time.Now(),
	}, nil
}
