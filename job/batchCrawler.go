package job

import (
	"digi-bot/crawler"
	"github.com/mileusna/crontab"
	"log"
)

func BatchCrawlerJob() {
	cTab := crontab.New()

	err := cTab.AddJob("0 0 * * *", crawler.StartBatchCrawl) // on 1st day of month
	if err != nil {
		log.Println(err)
		return
	}
}
