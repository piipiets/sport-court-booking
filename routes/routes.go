package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/piipiets/sport-court-booking/handler"
	"github.com/piipiets/sport-court-booking/middlewares"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	courtHandler *handler.CourtHandler,
	bookHandler *handler.BookingHandler,
) {
	router.POST("/login", userHandler.Login)
	router.POST("/sign-up", userHandler.SignUp)

	api := router.Group("/api/courts")
	api.Use(middlewares.JwtMiddleware())
	{
		api.POST("", courtHandler.Create)
		api.GET("", courtHandler.GetAll)
		api.GET("/:id", courtHandler.GetByID)
		api.PUT("/:id", courtHandler.Update)
		api.DELETE("/:id", courtHandler.Delete)
	}

	book := router.Group("/api/bookings")
	book.Use(middlewares.JwtMiddleware())
	{
		book.POST("", bookHandler.Create)
		book.GET("", bookHandler.GetAll)
		book.GET("/:id", bookHandler.GetByID)
		book.PUT("/status/:id", bookHandler.UpdateStatus)
	}
}
