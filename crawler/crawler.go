package crawler

import (
	"digi-bot/model"
	"encoding/json"
	"errors"
	"github.com/gocolly/colly"
)

// todo: create interface
func Crawl(url string) (model.Product, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("digikala.com", "www.digikala.com"),
	)

	var product model.Product
	c.OnHTML("script[type=\"application/ld+json\"]", func(e *colly.HTMLElement) {
		err := json.Unmarshal([]byte(e.Text), &product)
		if err != nil {
			println("err")
		}

		product.Url = url
		product.Price = product.Offer.Price
		product.OldPrice = product.Offer.OldPrice
	})

	err := c.Visit(url)
	if err != nil {
		println("visit error")
		return product, errors.New("visit error")
	}

	return product, nil

}
