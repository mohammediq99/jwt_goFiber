package database

import (
	"project1/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	connection, err := gorm.Open(mysql.Open("root:@/go_jwt"), &gorm.Config{})

	if err != nil {
		panic("Couldn't connect to DB")
	}
	DB = connection
	connection.AutoMigrate(&models.User{})
}
