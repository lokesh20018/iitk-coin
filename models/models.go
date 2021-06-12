// models/models.go

package models

import (
	"iitk-coin/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//defines the user in db
type User struct {
	gorm.Model
	Roll_no  string `json:"roll_no" gorm:"unique"`
	Password string `json:"password"`
}

// CreateUserRecord creates a user record in the database
func (user *User) CreateUserRecord() error {
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
