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
	digikalaCrawler = DigikalaCrawler{}
	torobCrawler    = TorobCrawler{}
)

func Crawl(url string) (model.ProductDto, error) {
	if strings.Contains(url, digikalaCrawler.GetDomain()) {
		return digikalaCrawler.Crawl(url)
	} else if strings.Contains(url, torobCrawler.GetDomain()) {
		return torobCrawler.Crawl(url)
	} else {
		return model.ProductDto{}, errors.New("crawler for that domain not found")
	}
}
