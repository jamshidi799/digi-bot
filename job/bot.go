package job

import (
	"digi-bot/bot/telegramBot"
	"sync"
)

func StartBot(group *sync.WaitGroup) {
	bot.InitTelegramBot(group)
}
