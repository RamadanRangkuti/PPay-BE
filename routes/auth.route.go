package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/controllers"
)

func AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
	}
}
