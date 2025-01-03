package db

import "myIris/models"

func Migrate() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Lead{})
	DB.AutoMigrate(&models.Customer{})
	DB.AutoMigrate(&models.Interaction{})
}
