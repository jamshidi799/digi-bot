package main

import (
	"digi-bot/job"
	"digi-bot/model/db"
	"digi-bot/service/bot"
	"sync"
)

func main() {
	db.Init()
	//go job.BatchCrawlerJob()

	var group sync.WaitGroup
	group.Add(1)

	go bot.StartBot(&group)

	group.Wait()

	group.Add(1)
	job.ChangeDetectorJob()
	group.Wait()
}
