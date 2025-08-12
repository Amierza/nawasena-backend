package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Admin(route *gin.Engine, adminHandler handler.IAdminHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/admins").Use(middleware.Authentication(jwtService), middleware.RouteAccessControl(jwtService))
	{
		routes.POST("/", adminHandler.Create)
		routes.GET("/", adminHandler.GetAll)
		routes.GET("/:id", adminHandler.GetDetail)
		routes.PATCH("/:id", adminHandler.Update)
		routes.DELETE("/:id", adminHandler.Delete)
	}
}
