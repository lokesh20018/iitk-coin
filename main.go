// main.go

package main

import (
	"github.com/lokesh20018/iitk-coin/controllers"

	"log"

	"github.com/lokesh20018/iitk-coin/models"

	"github.com/lokesh20018/iitk-coin/middlewares"

	"github.com/lokesh20018/iitk-coin/database"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	route := gin.Default()

	route.GET("/check", func(context *gin.Context) {
		context.String(200, "good to go")
	})

	route.POST("/login", controllers.Login)
	route.POST("/signup", controllers.Signup)
	route.POST("/init", middlewares.Authz_Admin(), controllers.Account_init)
	route.GET("/balance", middlewares.Authz(), controllers.GetBalance)
	route.POST("/transfer", middlewares.Authz(), controllers.Transfer)
	route.GET("/secretpage", middlewares.Authz(), controllers.Profile)

	// api_file := route.Group("/secretpage")
	// {
	// 	protected_route := api_file.Group("/").Use(middlewares.Authz())
	// 	{
	// 		protected_route.GET("/", controllers.Profile)
	// 	}
	// }

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

	err3 := database.InitDatabaseTrans()
	if err3 != nil {
		log.Fatalln("could not create Transfer database.... ", err3)
	}
	database.GlobalDBTrans.AutoMigrate(&models.Transaction{})

	route := setupRouter()
	route.Run(":8080")
}
