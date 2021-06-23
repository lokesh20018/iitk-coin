// database/database.go

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GlobalDB a global db object will be used across different packages
var GlobalDB *gorm.DB
var GlobalDBAcc *gorm.DB
var GlobalDBTrans *gorm.DB

// InitDatabase creates a sqlite db
func InitDatabase() (err error) {
	GlobalDB, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
	if err != nil {
		return
	}

	return
}

func InitDatabaseAcc() (err error) {
	GlobalDBAcc, err = gorm.Open(sqlite.Open("acc.db"), &gorm.Config{})
	if err != nil {
		return
	}
	return
}

func InitDatabaseTrans() (err error) {
	GlobalDBTrans, err = gorm.Open(sqlite.Open("trans.db"), &gorm.Config{})
	if err != nil {
		return
	}
	return
}
