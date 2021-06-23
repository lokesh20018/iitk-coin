// models/models.go

package models

import (
	"github.com/lokesh20018/iitk-coin/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//defines the user in db
type User struct {
	gorm.Model
	Roll_no  string `json:"roll_no" gorm:"unique"`
	Password string `json:"password"`
}

// the transfer model.... (to be used in transactions between accounts..)

// account struct...
type Account struct {
	gorm.Model
	Owner   string `json:"roll_no" gorm:"unique"`
	Balance int64  `json:"balance"`
}

// to record the entry...
type Transaction struct {
	gorm.Model
	FromAccountID string `json:"from_roll_no"`
	ToAccountID   string `json:"to_roll_no"`
	Amount        int64  `json:"amount"`
}

// CreateUserRecord creates a user record in the database
func (user *User) CreateUserRecord() error {
	//time.Sleep(8 * time.Second)
	result := database.GlobalDB.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// HashPassword encrypts user password
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	if err != nil {
		return err
	}
	println(bytes)
	user.Password = string(bytes)

	return nil
}

// CheckPassword checks user password
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}

// account init
func (account *Account) AccountInit() error {
	result := database.GlobalDBAcc.Create(&account)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (transaction *Transaction) TransactionRecord() error {
	result := database.GlobalDBTrans.Create(&transaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
