package main

import (
	"digi-bot/job"
	"digi-bot/model/db"
	"digi-bot/service/bot"
	"digi-bot/service/kafka"
	"sync"
)

func main() {
	db.Init()

	go kafka.InitProducer()
	defer kafka.FlushAndClose()

	var group sync.WaitGroup
	group.Add(1)

	go bot.StartBot(&group)

	group.Wait()

	group.Add(1)
	job.ChangeDetectorJob()
	group.Wait()

}
