package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GenerateSuccessResponse(ctx *gin.Context, message string) {
	ctx.JSON(
		http.StatusOK,
		GenerateSuccessMessage(message),
	)
}

func GenerateSuccessResponseWithData(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(
		http.StatusOK,
		GenerateSuccessMessageWithData(message, data),
	)
}

func GenerateErrorResponse(ctx *gin.Context, message string, httpStatusCode int) {
	ctx.AbortWithStatusJSON(
		httpStatusCode,
		GenerateErrorMessage(message),
	)
}

func GenerateSuccessMessage(message string) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    nil,
	}
}

func GenerateSuccessMessageWithData(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func GenerateErrorMessage(message string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	}
}
