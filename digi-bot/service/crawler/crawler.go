package crawler

import (
	"digi-bot/model"
	"errors"
	"strings"
)

type Crawler interface {
	Crawl(url string) (model.ProductDto, error)
	GetDomain() string
}

var (
	Crawlers = []Crawler{&DigikalaCrawler{}, &TorobCrawler{}}
)

func Crawl(url string) (model.ProductDto, error) {
	for _, crawler := range Crawlers {
		if strings.Contains(url, crawler.GetDomain()) {
			return crawler.Crawl(url)
		}
	}

	return model.ProductDto{}, errors.New("crawler for that domain not found")
}
