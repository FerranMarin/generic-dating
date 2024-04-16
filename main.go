package main

import (
	"github.com/FerranMarin/generic-dating/controllers"
	"github.com/FerranMarin/generic-dating/initializers"
	"github.com/FerranMarin/generic-dating/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
}

func main() {
	r := gin.Default()

	r.GET("/user/create", controllers.CreateRandom)
	r.POST("/user/login", controllers.Login)
	r.GET("/discover", middleware.RequireAuth, controllers.DiscoverUsers)
	r.POST("/swipe/:id", middleware.RequireAuth, controllers.DiscoverUsers)

	r.Run()
}
