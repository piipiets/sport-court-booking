package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/piipiets/sport-court-booking/handler"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
) {
	router.POST("/login", userHandler.Login)
}
