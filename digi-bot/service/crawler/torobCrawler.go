package crawler

import (
	"digi-bot/model"
	"github.com/gocolly/colly"
	"github.com/m7shapan/njson"
	"log"
)

type TorobCrawler struct {
}

func (TorobCrawler) Crawl(url string) (dto model.ProductDto, err error) {
	c := colly.NewCollector()

	var product torobProduct
	c.OnHTML(`script[type="application/ld+json"]`, func(e *colly.HTMLElement) {
		if e.Index == 0 {
			err = njson.Unmarshal([]byte(e.Text), &product)
			if err != nil {
				log.Println(err)
				return
			}
			dto.Url = url
			dto.Name = product.Name
			dto.Price = product.LowPrice / 10
			dto.Sku = product.Sku
			dto.Image = product.Image

			dto.Status = 1
			if dto.Price == 0 {
				dto.Status = 0
			}
		}
	})
	err = c.Visit(url)
	return
}

func (c *TorobCrawler) GetDomain() string {
	return "torob.com"
}

type torobProduct struct {
	ID       string
	Name     string `njson:"name"`
	Image    string `njson:"image"`
	LowPrice int    `njson:"offers.lowPrice"`
	Sku      string `njson:"sku"`
}
