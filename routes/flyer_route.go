package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Flyer(route *gin.Engine, flyerHandler handler.IFlyerHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/flyers")
	{
		routes.GET("", flyerHandler.GetAll)
		routes.GET("/:id", flyerHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("", flyerHandler.Create)
			routes.PATCH("/:id", flyerHandler.Update)
			routes.DELETE("/:id", flyerHandler.Delete)
		}
	}
}
