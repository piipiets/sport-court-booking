package handler

import (
	"fmt"

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

func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Bind error:", err)
		message := common.GetValidationError(err)
		common.GenerateErrorResponse(c, message)
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			common.GenerateErrorResponse(c, "invalid email or password")
			return
		}

		common.GenerateErrorResponse(c, "failed to login")
		return
	}

	common.GenerateSuccessResponseWithData(c, "login successful", result)
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var req request.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Bind error:", err)
		message := common.GetValidationError(err)
		common.GenerateErrorResponse(c, message)
		return
	}

	err := h.service.SignUp(c.Request.Context(), req)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error())
		return
	}

	common.GenerateSuccessResponse(c, "login successful")
}
