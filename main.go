package main

import (
	"log"
	"os"

	"github.com/Amierza/nawasena-backend/cmd"
	"github.com/Amierza/nawasena-backend/config/database"
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/Amierza/nawasena-backend/repository"
	"github.com/Amierza/nawasena-backend/routes"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	var (
		jwt = jwt.NewJWT()

		// Auth
		authRepo    = repository.NewAuthRepository(db)
		authService = service.NewAuthService(authRepo, jwt)
		authHandler = handler.NewAuthHandler(authService)

		// Files
		fileService = service.NewFileService()
		fileHandler = handler.NewFileHandler(fileService)

		// Admin
		adminRepo    = repository.NewAdminRepository(db)
		adminService = service.NewAdminService(adminRepo, jwt)
		adminHandler = handler.NewAdminHandler(adminService)

		// Position
		positionRepo    = repository.NewPositionRepository(db)
		positionService = service.NewPositionService(positionRepo, jwt)
		positionHandler = handler.NewPositionHandler(positionService)

		// Member
		memberRepo    = repository.NewMemberRepository(db)
		memberService = service.NewMemberService(memberRepo, jwt)
		memberHandler = handler.NewMemberHandler(memberService)

		// Achievement
		achievementRepo    = repository.NewAchievementRepository(db)
		achievementService = service.NewAchievementService(achievementRepo, jwt)
		achievementHandler = handler.NewAchievementHandler(achievementService)

		// Ship
		shipRepo    = repository.NewShipRepository(db)
		shipService = service.NewShipService(shipRepo, jwt)
		shipHandler = handler.NewShipHandler(shipService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.Auth(server, authHandler, jwt)
	routes.File(server, fileHandler, jwt)
	routes.Admin(server, adminHandler, jwt)
	routes.Position(server, positionHandler, jwt)
	routes.Member(server, memberHandler, jwt)
	routes.Achievement(server, achievementHandler, jwt)
	routes.Ship(server, shipHandler, jwt)

	server.Static("/assets", "./assets")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
