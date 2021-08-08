package crawler

import (
	"digi-bot/messageCreator"
	"digi-bot/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
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
			log.Println("err")
		}
		product.Url = url
	})

	c.OnHTML(".c-product", func(e *colly.HTMLElement) {
		price := e.ChildText(".c-product__seller-price-pure.js-price-value")
		product.Price = messageCreator.FixNumber(price)

		oldPrice := e.ChildText(".c-product__seller-price-prev.js-rrp-price")
		product.OldPrice = messageCreator.FixNumber(oldPrice)

		fmt.Printf("%s %s\n", price, oldPrice)

		desc1 := e.ChildText(".c-product__user-suggestion-line")
		desc1 = messageCreator.CleaningString(desc1)
		desc1 = strings.Split(desc1, ".")[0]
		product.Desc1 = messageCreator.CleaningString(desc1)

		fmt.Println(e.ChildText("c-product__user-suggestion-line"))
		desc2 := e.ChildText(".c-product__engagement-rating")
		product.Desc2 = "امتیاز " + messageCreator.CleaningString(desc2)
	})

	err := c.Visit(url)
	if err != nil {
		println("visit error")
		return product, errors.New("visit error")
	}

	return product, nil

}
