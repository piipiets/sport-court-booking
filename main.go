package main

import (
	"github.com/gin-gonic/gin"

	"github.com/piipiets/sport-court-booking/configs"
	"github.com/piipiets/sport-court-booking/databases/connection"
	"github.com/piipiets/sport-court-booking/databases/migration"
	"github.com/piipiets/sport-court-booking/handler"
	"github.com/piipiets/sport-court-booking/repository"
	"github.com/piipiets/sport-court-booking/routes"
	"github.com/piipiets/sport-court-booking/service"
)

func main() {
	configs.Initiator()

	connection.Initiator()
	defer connection.DBConnections.Close()

	conn := connection.DBConnections
	migration.Initiator(conn)

	// repository
	userRepository := repository.NewUserRepository(conn)
	courtRepository := repository.NewCourtRepository(conn)

	// service
	userService := service.NewUserService(userRepository)
	courtService := service.NewCourtService(courtRepository)

	// handler
	userHandler := handler.NewUserHandler(userService)
	courtHandler := handler.NewCourtHandler(courtService)

	// router
	router := gin.Default()

	// routes
	routes.SetupRoutes(
		router,
		userHandler,
		courtHandler,
	)

	router.Run(":8080")
}
