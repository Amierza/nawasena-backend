package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Partner(route *gin.Engine, partnerHandler handler.IPartnerHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/partners")
	{
		routes.GET("", partnerHandler.GetAll)
		routes.GET("/:id", partnerHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("", partnerHandler.Create)
			routes.PATCH("/:id", partnerHandler.Update)
			routes.DELETE("/:id", partnerHandler.Delete)
		}
	}
}
