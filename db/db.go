package db

import (
	"digi-bot/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init() *gorm.DB {
	dsn := "root:mohammad79@tcp(127.0.0.1:3306)/digiBot?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	migration()
	return DB
}

func migration() {
	_ = DB.AutoMigrate(&model.UserModel{})
	_ = DB.AutoMigrate(&model.ObjectModel{})
}

func GetAllProduct() []model.ObjectModel {
	var objects []model.ObjectModel
	DB.Find(&objects)
	return objects
}

func GetUserById(userId int) model.UserModel {
	var user model.UserModel
	DB.Where(model.UserModel{ID: userId}).First(&user)
	return user
}
