package db

import (
	"digi-bot/model"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func Init() *gorm.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DbPassword := os.Getenv("DB_PASSWORD")
	DbHost := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("root:%s@tcp(%s:3306)/digiBot?charset=utf8mb4&parseTime=True&loc=Local", DbPassword, DbHost)
	println(dsn)
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
	_ = DB.AutoMigrate(&model.ProductModel{})
}
