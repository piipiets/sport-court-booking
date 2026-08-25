package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piipiets/sport-court-booking/handler"
	"github.com/piipiets/sport-court-booking/middlewares"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	courtHandler *handler.CourtHandler,
	bookHandler *handler.BookingHandler,
	paymentHandler *handler.PaymentHandler,
) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Sport Court Booking API",
			"status":  "running",
			"version": "1.0.0",
			"endpoints": gin.H{
				"auth": gin.H{
					"login":  "POST /login",
					"signup": "POST /sign-up",
				},
				"courts":   "/api/courts (with JWT)",
				"bookings": "/api/bookings (with JWT)",
				"payments": "/api/payments (with JWT)",
			},
			"timestamp": time.Now().Unix(),
		})
	})

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

	payment := router.Group("/api/payments")
	payment.Use(middlewares.JwtMiddleware())
	{
		payment.POST("", paymentHandler.Create)
		payment.GET("/:booking_id", paymentHandler.GetByBookingID)
		payment.GET("", paymentHandler.GetAllPaymentByUserId)
	}
}
