package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/piipiets/sport-court-booking/helpers/common"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/service"
)

type UserHandler struct {
	service service.Service
}

func NewUserHandler(service service.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Bind error:", err)
		common.GenerateErrorResponse(c, "invalid request body")
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
