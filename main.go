package main

import (
	"digi-bot/bot"
	"digi-bot/config"
)

func main() {
	config.Init()
	bot.Run()
	//crawler.Crawl()
}
