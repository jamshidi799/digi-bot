package crawler

import (
	"digi-bot/model"
	"encoding/json"
	"errors"
	"github.com/gocolly/colly"
)

func Crawl(url string) (model.Object, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("digikala.com", "www.digikala.com"),
	)

	var obj model.Object
	c.OnHTML("script[type=\"application/ld+json\"]", func(e *colly.HTMLElement) {
		err := json.Unmarshal([]byte(e.Text), &obj)
		if err != nil {
			println("err")
		}

		obj.Url = url
		obj.Price = obj.Offer.Price
		obj.OldPrice = obj.Offer.OldPrice
	})

	err := c.Visit(url)
	if err != nil {
		println("visit error")
		return obj, errors.New("visit error")
	}

	return obj, nil

}
