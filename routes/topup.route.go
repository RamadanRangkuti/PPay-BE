package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/controllers"
)

func TopupRoute(router *gin.Engine) {
	topupGroup := router.Group("/topup")
	{
		topupGroup.POST("", controllers.Topup)
	}
}
