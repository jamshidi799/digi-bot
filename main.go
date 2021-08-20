package main

import (
	"digi-bot/bot"
	"digi-bot/db"
	"digi-bot/job"
	"sync"
)

func main() {
	db.Init()
	go job.BatchCrawlerJob()

	var group sync.WaitGroup
	group.Add(1)

	go bot.Run(&group)

	group.Wait()

	group.Add(1)
	job.ChangeDetectorJob()
	group.Wait()
}
