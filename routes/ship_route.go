package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Ship(route *gin.Engine, shipHandler handler.IShipHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/ships")
	{
		routes.GET("", shipHandler.GetAll)
		routes.GET("/:id", shipHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("", shipHandler.Create)
			routes.PATCH("/:id", shipHandler.Update)
			routes.DELETE("/:id", shipHandler.Delete)
		}
	}
}
