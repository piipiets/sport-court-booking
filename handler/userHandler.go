package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piipiets/sport-court-booking/helpers/common"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// @Summary      User login
// @Description  Authenticate a user and return a JWT access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.LoginRequest true "Login credentials"
// @Success      200 {object} common.APIResponse{data=response.LoginResponse}
// @Failure      400 {object} common.APIResponse
// @Failure      401 {object} common.APIResponse
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Bind error:", err)
		message := common.GetValidationError(err)
		common.GenerateErrorResponse(c, message, http.StatusBadRequest)
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			common.GenerateErrorResponse(c, "invalid email or password", http.StatusUnauthorized)
			return
		}

		common.GenerateErrorResponse(c, "failed to login", http.StatusUnauthorized)
		return
	}

	common.GenerateSuccessResponseWithData(c, "login successful", result)
}

// @Summary      User sign up
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.SignUpRequest true "Sign up payload"
// @Success      200 {object} common.APIResponse
// @Failure      400 {object} common.APIResponse
// @Router       /sign-up [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var req request.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Bind error:", err)
		message := common.GetValidationError(err)
		common.GenerateErrorResponse(c, message, http.StatusBadRequest)
		return
	}

	err := h.service.SignUp(c.Request.Context(), req)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	common.GenerateSuccessResponse(c, "login successful")
}
