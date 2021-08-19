package user

import (
	"digi-bot/bot"
	"digi-bot/db"
	"log"
)

func SendProductUpdateToUsers(productId int, message string, changeLevel int) {
	usersId := db.GetAllUsersIdByProductId(productId, changeLevel)

	log.Printf("user affected: %d", len(usersId))

	bot.SendUpdateForUsers(usersId, productId, message)
}
