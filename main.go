// main.go

package main

import (
	"iitk-coin/controllers"
	"iitk-coin/database"
	"iitk-coin/models"
	"log"

	"iitk-coin/middlewares"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	route := gin.Default()

	route.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	route.POST("/login", controllers.Login)
	route.POST("/signup", controllers.Signup)

	api_file := route.Group("/secretpage")
	{
		protected_route := api_file.Group("/").Use(middlewares.Authz())
		{
			protected_route.GET("/", controllers.Profile)
		}
	}

	return route
}

func main() {
	err := database.InitDatabase()
	if err != nil {
		log.Fatalln("could not create database", err)
	}

	database.GlobalDB.AutoMigrate(&models.User{})

	route := setupRouter()
	route.Run(":8080")
}
