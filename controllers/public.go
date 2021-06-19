// controllers/public.go

package controllers

import (
	"iitk-coin/auth"
	"iitk-coin/database"
	"iitk-coin/models"
	"log"

	"github.com/gin-gonic/gin"
	//"github.com/lokesh20018/iitk-coin/database"
	"gorm.io/gorm"
)

// LoginPayload
type LoginPayload struct {
	Roll_no  string `json:"roll_no"`
	Password string `json:"password"`
}

// LoginResponse token
type LoginResponse struct {
	Token string `json:"token"`
}

// add money / init payload
type InitPayload struct {
	Owner   string `json:"roll_no"`
	Balance int64  `json:"balance"`
}

// query balance paylode..
type BalancePayload struct {
	Owner string `json:"roll_no"`
}

type BalanceResponse struct {
	Owner   string `json:"roll_no"`
	Balance int64  `json:"balance`
}

// creates a user in db
func Signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)

		context.JSON(400, gin.H{
			"msg": "invalid json",
		})
		context.Abort()

		return
	}

	err = user.HashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())

		context.JSON(500, gin.H{
			"msg": "error hashing password",
		})
		context.Abort()

		return
	}

	err = user.CreateUserRecord()
	if err != nil {
		log.Println(err)

		context.JSON(500, gin.H{
			"msg": "error creating user",
		})
		context.Abort()

		return
	}

	context.JSON(200, user)
}

// logs users in
func Login(context *gin.Context) {
	var payload LoginPayload
	var user models.User

	err := context.ShouldBindJSON(&payload)
	if err != nil {
		context.JSON(400, gin.H{
			"msg": "invalid json",
		})
		context.Abort()
		return
	}

	result := database.GlobalDB.Where("Roll_no = ?", payload.Roll_no).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		context.JSON(401, gin.H{
			"msg": "invalid user credentials",
		})
		context.Abort()
		return
	}

	err = user.CheckPassword(payload.Password)
	if err != nil {
		log.Println(err)
		context.JSON(401, gin.H{
			"msg": "invalid user credentials",
		})
		context.Abort()
		return
	}

	jwtWrapper := auth.JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 1,
	}

	signedToken, err := jwtWrapper.GenerateToken(user.Roll_no)
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{
			"msg": "error signing token",
		})
		context.Abort()
		return
	}

	tokenResponse := LoginResponse{
		Token: signedToken,
	}

	context.JSON(200, tokenResponse)

	return
}

//Account INIT
func Account_init(context *gin.Context) {
	//context2 := context

	var payload InitPayload
	var account models.Account
	//var account2 models.Account

	err := context.ShouldBindJSON(&payload)
	if err != nil {
		context.JSON(400, gin.H{
			"msg": "invalid json",
		})
		context.Abort()
		return
	}

	result := database.GlobalDBAcc.Where("owner = ?", payload.Owner).First(&account)

	if result.Error == gorm.ErrRecordNotFound {
		account.Owner = payload.Owner
		account.Balance = payload.Balance
		err = account.AccountInit()
		if err != nil {
			log.Println(err)
			context.JSON(500, gin.H{
				"msg": "error creating user",
			})
			context.Abort()

			return
		}
		//context.JSON(200, account)

	} else {
		account.Balance += payload.Balance
		database.GlobalDBAcc.Save(&account)
		//context.JSON(200, account)

	}

	return
}

func GetBalance(context *gin.Context) {
	var payload BalancePayload
	var account models.Account

	err := context.ShouldBindJSON(&payload)
	if err != nil {
		context.JSON(400, gin.H{
			"msg": "invalid json",
		})
		context.Abort()
		return
	}

	result := database.GlobalDBAcc.Where("owner = ?", payload.Owner).First(&account)

	if result.Error == gorm.ErrRecordNotFound {
		context.JSON(401, gin.H{
			"msg": "invalid user credentials",
		})
		context.Abort()
		return
	}
	var response BalanceResponse
	response.Owner = account.Owner
	response.Balance = account.Balance

	context.JSON(200, response)
	return
}
