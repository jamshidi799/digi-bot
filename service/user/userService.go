package user

import (
	"digi-bot/bot/telegramBot"
	"digi-bot/db"
	"log"
)

func SendProductUpdateToUsers(productId int, message string, changeLevel int) {
	usersId := db.GetAllUsersIdByProductId(productId, changeLevel)

	log.Printf("user affected: %d", len(usersId))

	telegramBot.GetBot().SendUpdateForUsers(usersId, productId, message)
}
