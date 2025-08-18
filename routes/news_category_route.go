package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func NewsCategory(route *gin.Engine, newsCategoryHandler handler.INewsCategoryHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/news-categories")
	{
		routes.GET("", newsCategoryHandler.GetAll)
		routes.GET("/:id", newsCategoryHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("", newsCategoryHandler.Create)
			routes.PATCH("/:id", newsCategoryHandler.Update)
			routes.DELETE("/:id", newsCategoryHandler.Delete)
		}
	}
}
