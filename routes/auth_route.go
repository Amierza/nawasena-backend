package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/gin-gonic/gin"
)

func Auth(route *gin.Engine, authHandler handler.IAuthHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1")
	{
		routes.POST("/login", authHandler.Login)
		routes.POST("/refresh-token", authHandler.RefreshToken)
	}
}
