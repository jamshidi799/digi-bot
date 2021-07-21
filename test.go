package main

import (
	"digi-bot/config"
	"digi-bot/model"
)

func main() {
	db := config.Init()
	test := model.TestModel{Username: "mohammad", FirstName: "jamshidi", LastName: "amoo mama"}
	db.Create(&test)

	//defer db.
}
