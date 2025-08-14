package routes

import (
	"github.com/Amierza/nawasena-backend/handler"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/middleware"
	"github.com/gin-gonic/gin"
)

func File(route *gin.Engine, fileHandler handler.IFileHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/uploads", middleware.Authentication(jwtService))
	{
		routes.POST("", fileHandler.UploadFiles)
	}
}
