package crawler

import "digi-bot/model"

type Crawler interface {
	Crawl(url string) (model.ProductDto, error)
	BatchCrawl()
}
