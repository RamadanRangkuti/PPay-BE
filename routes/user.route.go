package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/controllers"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("", controllers.CreateUser)
		userGroup.GET("", controllers.GetUsers)
		userGroup.GET("/:id", controllers.GetUserByID)
		userGroup.PATCH("/:id", controllers.UpdateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}
}
