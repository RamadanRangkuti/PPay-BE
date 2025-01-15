package main

import (
	"github.com/ppay/initializers"
	"github.com/ppay/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDb()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
}
