package main

import (
	"log"
	"os"

	"github.com/Amierza/nawasena-backend/cmd"
	"github.com/Amierza/nawasena-backend/config/database"
	_ "github.com/Amierza/nawasena-backend/docs"
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/Amierza/nawasena-backend/repository"
	"github.com/Amierza/nawasena-backend/routes"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Nawasena API
// @version         1.0
// @description     API documentation for Nawasena project
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      nawasena-backend-production.up.railway.app
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	var (
		jwtService = jwt.NewJWTService()

		// Auth
		authRepo    = repository.NewAuthRepository(db)
		authService = service.NewAuthService(authRepo, jwtService)
		authHandler = handler.NewAuthHandler(authService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.Auth(server, authHandler, jwtService)
	// swagger endpoint
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
