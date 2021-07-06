// controllers/public.go

package controllers

import (
	"log"

	"github.com/lokesh20018/iitk-coin/auth"
	"github.com/lokesh20018/iitk-coin/database"
	"github.com/lokesh20018/iitk-coin/models"

	"github.com/gin-gonic/gin"
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

// Transfer payload..
type TransferPayload struct {
	FromAccountID string `json:"from_roll_no"`
	ToAccountID   string `json:"to_roll_no"`
	// must be positive
	Amount int64 `json:"amount"`
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

	var account models.Account
	account.Owner = user.Roll_no
	account.Balance = 0
	err = account.AccountInit()
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{
			"msg": "error creating user account",
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
	var payload InitPayload
	var account models.Account
	var transaction models.Transaction
	err := context.ShouldBindJSON(&payload)
	if err != nil {
		context.JSON(400, gin.H{
			"msg": "invalid json",
		})
		context.Abort()
		return
	}
	if payload.Balance > 30000 {
		context.JSON(400, gin.H{
			"msg": "Upper limit to one time transaction is 30000 coins",
		})
		context.Abort()
		return
	}
	result := database.GlobalDBAcc.Where("owner = ?", payload.Owner).First(&account)

	if result.Error == gorm.ErrRecordNotFound {
		context.JSON(500, gin.H{
			"msg": "Account not found",
		})
	} else if account.Balance+payload.Balance > 1000000 {
		context.JSON(400, gin.H{
			"msg": "Account upper limit reached ",
		})
		context.Abort()
		return
	} else {
		account.Balance += payload.Balance
		database.GlobalDBAcc.Save(&account)
		//context.JSON(200, account)
		context.JSON(200, gin.H{
			"msg": "added amount",
		})
	}
	transaction.Amount = payload.Balance
	transaction.FromAccountID = "Admin_reward"
	transaction.ToAccountID = payload.Owner
	transaction.TransactionRecord()

	return
}

// read balance..
func GetBalance(context *gin.Context) {
	Roll_no, _ := context.Get("roll_no")
	// context.JSON(401, gin.H{
	// 	"msg": Roll_no,
	// })
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
	if Roll_no != payload.Owner {
		context.JSON(401, gin.H{
			"msg": "not authorised to view balance",
		})
		context.Abort()
		return
	}
	//time.Sleep(5 * time.Second)
	result := database.GlobalDBAcc.Where("owner = ?", payload.Owner).First(&account)
	//time.Sleep(1 * time.Second)

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

// trasaction

func Transfer(context *gin.Context) {
	Roll_no, _ := context.Get("roll_no")
	var payload TransferPayload
	var FromAcc models.Account
	var ToAcc models.Account
	var transaction models.Transaction

	err := context.ShouldBindJSON(&payload)
	if err != nil {
		context.JSON(400, gin.H{
			"msg": "invalid json",
		})
		context.Abort()
		return
	}

	if Roll_no != payload.FromAccountID {
		context.JSON(401, gin.H{
			"msg": "not authorised to transfer",
		})
		context.Abort()
		return
	}

	if payload.Amount > 10000 {
		context.JSON(400, gin.H{
			"msg": "Upper limit to one time transaction is 10000 coins",
		})
		context.Abort()
		return
	}
	if payload.Amount < 100 {
		context.JSON(400, gin.H{
			"msg": "lower limit of transaction is 100 coins",
		})
		context.Abort()
		return
	}

	database.GlobalDBAcc.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("owner = ?", payload.FromAccountID).First(&FromAcc)

		if result.Error == gorm.ErrRecordNotFound {
			context.JSON(500, gin.H{
				"msg": "error finding Sender account",
			})
			context.Abort()
			return result.Error
			//context.JSON(200, account)
		}
		if FromAcc.Balance < payload.Amount {

			context.JSON(500, gin.H{
				"msg": "account balance low",
			})
			context.Abort()
			tx.Rollback()
			return nil
		}
		FromAcc.Balance -= payload.Amount
		tx.Save(&FromAcc)
		//time.Sleep(8 * time.Second)
		result = tx.Where("owner = ?", payload.ToAccountID).First(&ToAcc)

		if result.Error == gorm.ErrRecordNotFound {
			context.JSON(500, gin.H{
				"msg": "error finding reciever account",
			})
			context.Abort()
			return result.Error
			//context.JSON(200, account)
		}
		txn_charge := 10
		if payload.Amount > 500 {
			txn_charge = (int(payload.Amount) / 100) * 2
		}
		payload.Amount -= int64(txn_charge)

		if ToAcc.Balance+payload.Amount > 1000000 {
			context.JSON(400, gin.H{
				"msg": "Account upper limit reached",
			})
			context.Abort()
			tx.Rollback()
			return nil
		}
		ToAcc.Balance += payload.Amount

		context.JSON(200, gin.H{
			"msg":              "transfer successful ",
			"Amount initiated": payload.Amount + int64(txn_charge),
			"txn charge":       txn_charge,
			"Amount recieved":  payload.Amount,
		})
		// return nil will commit the whole transaction
		transaction.Amount = payload.Amount
		transaction.FromAccountID = payload.FromAccountID
		transaction.ToAccountID = payload.ToAccountID
		transaction.TransactionRecord()
		tx.Save(&ToAcc)
		return nil
	})
}
