package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func News(route *gin.Engine, newsHandler handler.INewsHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/news")
	{
		routes.GET("/", newsHandler.GetAll)
		routes.GET("/:id", newsHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("/", newsHandler.Create)
			routes.PATCH("/:id", newsHandler.Update)
			routes.DELETE("/:id", newsHandler.Delete)
		}
	}
}
