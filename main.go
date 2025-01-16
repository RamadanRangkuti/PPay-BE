package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/initializers"
	"github.com/ppay/routes"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	routes.UserRoutes(r)
	routes.TopupRoute(r)
	routes.AuthRoutes(r)
	r.Run()
}
