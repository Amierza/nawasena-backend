package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Member(route *gin.Engine, memberHandler handler.IMemberHandler, jwt jwt.IJWT) {
	routes := route.Group("/api/v1/members")
	{
		routes.GET("/", memberHandler.GetAll)
		routes.GET("/:id", memberHandler.GetDetail)

		routes.Use(middleware.Authentication(jwt), middleware.RouteAccessControl(jwt))
		{
			routes.POST("/", memberHandler.Create)
			routes.PATCH("/:id", memberHandler.Update)
			routes.DELETE("/:id", memberHandler.Delete)
		}
	}
}
