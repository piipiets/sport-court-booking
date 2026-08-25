package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/piipiets/sport-court-booking/configs"
	"github.com/piipiets/sport-court-booking/databases/connection"
	"github.com/piipiets/sport-court-booking/databases/migration"
	apphandler "github.com/piipiets/sport-court-booking/handler"
	"github.com/piipiets/sport-court-booking/repository"
	"github.com/piipiets/sport-court-booking/routes"
	"github.com/piipiets/sport-court-booking/service"
)

var (
	router    *gin.Engine
	initOnce  sync.Once
	initPanic any
)

func Handler(writer http.ResponseWriter, request *http.Request) {
	initOnce.Do(func() {
		defer func() {
			initPanic = recover()
		}()
		router = buildRouter()
	})

	if initPanic != nil {
		panic(initPanic)
	}

	router.ServeHTTP(writer, request)
}

func buildRouter() *gin.Engine {
	configs.Initiator()
	connection.Initiator()

	conn := connection.DBConnections
	migration.Initiator(conn)

	userRepository := repository.NewUserRepository(conn)
	courtRepository := repository.NewCourtRepository(conn)
	bookingRepository := repository.NewBookingRepository(conn)
	paymentRepository := repository.NewPaymentRepository(conn)

	userService := service.NewUserService(userRepository)
	courtService := service.NewCourtService(courtRepository)
	bookingService := service.NewBookingService(bookingRepository, courtRepository)
	paymentService := service.NewPaymentService(paymentRepository, bookingRepository)

	userHandler := apphandler.NewUserHandler(userService)
	courtHandler := apphandler.NewCourtHandler(courtService)
	bookingHandler := apphandler.NewBookingHandler(bookingService)
	paymentHandler := apphandler.NewPaymentHandler(paymentService)

	ginRouter := gin.Default()
	routes.SetupRoutes(
		ginRouter,
		userHandler,
		courtHandler,
		bookingHandler,
		paymentHandler,
	)

	return ginRouter
}
