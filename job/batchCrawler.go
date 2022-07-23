package job

import (
	"digi-bot/service/crawler/digikalaCrawler"
	"github.com/mileusna/crontab"
	"log"
)

func BatchCrawlerJob() {
	cTab := crontab.New()

	err := cTab.AddJob("0 0 */8 * *", crawler.DigikalaCrawler{}.BatchCrawl) // on 1st day of month
	if err != nil {
		log.Println(err)
		return
	}
}
