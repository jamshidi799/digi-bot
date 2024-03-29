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

var database *gorm.DB

func Init() *gorm.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	DbUsername := os.Getenv("DB_USERNAME")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbHost := os.Getenv("DB_HOST")
	DbName := os.Getenv("DB_NAME")
	DbTlsEnabled := os.Getenv("DB_TLS_ENABLED")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s", DbUsername, DbPassword, DbHost, DbName, DbTlsEnabled)
	//println(dsn)

	for i := 1; i <= 5; i++ {
		log.Printf("connect to db try %d", i)
		database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
	return database
}

func migration() {
	_ = database.AutoMigrate(&model.Pivot{})
	_ = database.AutoMigrate(&model.History{})
	_ = database.AutoMigrate(&model.User{})
	_ = database.AutoMigrate(&model.Product{})
}
