package models

import (
	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/go_api"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Users{}, &Photo{})

	DB = db

}
