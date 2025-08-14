package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Position(route *gin.Engine, positionHandler handler.IPositionHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/positions")
	{
		routes.GET("", positionHandler.GetAll)
		routes.GET("/:id", positionHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("", positionHandler.Create)
			routes.PATCH("/:id", positionHandler.Update)
			routes.DELETE("/:id", positionHandler.Delete)
		}
	}
}
