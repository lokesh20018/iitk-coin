// main.go

package main

import (
	"iitk-coin/controllers"

	"iitk-coin/models"
	"log"

	"iitk-coin/middlewares"

	"iitk-coin/database"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	route := gin.Default()

	route.GET("/check", func(context *gin.Context) {
		context.String(200, "good to go")
	})

	route.POST("/login", controllers.Login)
	route.POST("/signup", controllers.Signup)
	route.POST("/init", controllers.Account_init)
	route.GET("/balance", controllers.GetBalance)
	route.POST("/transfer", controllers.Transfer)
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

	err2 := database.InitDatabaseAcc()
	if err2 != nil {
		log.Fatalln("could not create Acc ", err2)
	}
	database.GlobalDBAcc.AutoMigrate(&models.Account{})

	route := setupRouter()
	route.Run(":8080")
}
