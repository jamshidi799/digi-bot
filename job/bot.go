package job

import (
	bot "digi-bot/bot/telegramBot"
	"sync"
)

func StartBot(group *sync.WaitGroup) {
	var tlBot bot.TelegramBot
	tlBot.Init(group)
}
