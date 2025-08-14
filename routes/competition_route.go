package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Competition(route *gin.Engine, competitionHandler handler.ICompetitionHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/competitions")
	{
		routes.GET("", competitionHandler.GetAll)
		routes.GET("/:id", competitionHandler.GetDetail)

		routes.Use(middleware.Authentication(jwtService))
		{
			routes.POST("", competitionHandler.Create)
			routes.PATCH("/:id", competitionHandler.Update)
			routes.DELETE("/:id", competitionHandler.Delete)
		}
	}
}
