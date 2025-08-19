package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AchievementCategory(route *gin.Engine, achievementCategoryHandler handler.IAchievementCategoryHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/achievement-categories")
	{
		routes.GET("", achievementCategoryHandler.GetAll)
		routes.GET("/:id", achievementCategoryHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("", achievementCategoryHandler.Create)
			routes.PATCH("/:id", achievementCategoryHandler.Update)
			routes.DELETE("/:id", achievementCategoryHandler.Delete)
		}
	}
}
