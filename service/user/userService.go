package user

import (
	"digi-bot/bot/telegramBot"
	"digi-bot/db"
	"log"
)

func SendProductUpdateToUsers(productId int, message string, changeLevel int, available bool) {
	usersId := db.GetAllUsersIdByProductId(productId, changeLevel)

	log.Printf("user affected: %d", len(usersId))

	bot.GetTelegramBot().SendUpdateForUsers(usersId, productId, message, available)
}
