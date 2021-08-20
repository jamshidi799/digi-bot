package main

import (
	"digi-bot/bot"
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/job"
	"log"
	"sync"
)

func main() {
	db.Init()
	var group sync.WaitGroup
	group.Add(1)
	go bot.Run(&group)
	group.Wait()
	log.Println("bot started")

	//go crawler.BatchCrawlScheduler()
	go crawler.StartBatchCrawl()
	group.Add(1)
	job.ChangeDetector()
	group.Wait()
}
