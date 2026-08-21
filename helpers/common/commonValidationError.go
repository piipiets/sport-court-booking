package common

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetValidationError(err error) string {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		e := validationErrors[0]

		field := e.Field()
		tag := e.Tag()

		switch tag {
		case "required":
			return strings.ToLower(field) + " is required"
		case "email":
			return strings.ToLower(field) + " must be a valid email"
		case "min":
			return strings.ToLower(field) + " must be at least " + e.Param() + " characters"
		case "eqfield":
			return strings.ToLower(field) + " must match " + e.Param()
		}

		return strings.ToLower(field) + " is invalid"
	}

	return "invalid request body"
}
