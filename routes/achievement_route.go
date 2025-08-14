package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Achievement(route *gin.Engine, achievementHandler handler.IAchievementHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/achievements")
	{
		routes.GET("", achievementHandler.GetAll)
		routes.GET("/:id", achievementHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("", achievementHandler.Create)
			routes.PATCH("/:id", achievementHandler.Update)
			routes.DELETE("/:id", achievementHandler.Delete)
		}
	}
}
