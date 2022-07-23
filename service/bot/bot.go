package bot

import (
	"digi-bot/service/bot/telegramBot"
	"sync"
)

type Bot interface {
	SendUpdateForUsers(usersId []int, productId int, message string, available bool)
}

func StartBot(group *sync.WaitGroup) {
	bot.InitTelegramBot(group)
}
