// controllers/protected.go

package controllers

import (
	"iitk-coin/database"
	"iitk-coin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Profile returns dummy user data after auth
func Profile(context *gin.Context) {
	var user models.User

	Roll_no, _ := context.Get("roll_no") //  authorization middleware

	response := database.GlobalDB.Where("roll_no = ?", Roll_no.(string)).First(&user)

	if response.Error == gorm.ErrRecordNotFound {
		context.JSON(404, gin.H{
			"msg": "user not found",
		})
		context.Abort()
		return
	}

	if response.Error != nil {
		context.JSON(500, gin.H{
			"msg": "could not get user profile",
		})
		context.Abort()
		return
	}

	user.Password = "Hidden for security purpose..."

	context.JSON(200, user)

	return
}
