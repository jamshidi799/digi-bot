package db

import (
	"digi-bot/model"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func Init() *gorm.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	DbPassword := os.Getenv("DB_PASSWORD")
	DbHost := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("root:%s@tcp(%s:3306)/digiBot?charset=utf8mb4&parseTime=True&loc=Local", DbPassword, DbHost)
	//println(dsn)

	for i := 1; i <= 5; i++ {
		log.Printf("connect to db try %d", i)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println(err)
		} else {
			log.Println("connected to mysql")
			break
		}

		time.Sleep(time.Second * 10)
	}

	if err != nil {
		log.Fatal(err)
	}

	migration()
	return DB
}

func migration() {
	_ = DB.AutoMigrate(&model.User{})
	_ = DB.AutoMigrate(&model.Product{})
	_ = DB.AutoMigrate(&model.Pivot{})
	_ = DB.AutoMigrate(&model.History{})
	_ = DB.AutoMigrate(&model.BulkHistory{})
}
