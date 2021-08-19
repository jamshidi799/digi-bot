package crawler

import (
	"digi-bot/messageCreator"
	"digi-bot/model"
	"encoding/json"
	"errors"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

// todo: create interface
func Crawl(url string) (model.ProductDto, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("digikala.com", "www.digikala.com"),
	)

	var product model.ProductDto
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

		product.Status = 1
		if product.Price == 0 {
			product.Status = 0
		}

		//fmt.Printf("%s %s\n", price, oldPrice)

		desc1 := e.ChildText(`div[class="c-product__headline--gallery "]`)
		product.Desc1 = messageCreator.CleaningString(desc1)

		if desc1 == "" {
			seconds := e.ChildAttr(`.c-product-gallery__timer.js-counter`, "data-countdownseconds")
			if seconds != "" {
				product.Desc1 = messageCreator.CreateAmazingOfferText(seconds)
				product.Status = 2
			}
		}

		desc2 := e.ChildText(".c-product__user-suggestion-line")
		desc2 = messageCreator.CleaningString(desc2)
		desc2 = strings.Split(desc2, ".")[0]
		product.Desc2 = messageCreator.CleaningString(desc2)

		desc3 := e.ChildText(".c-product__engagement-rating")
		if desc3 == "" {
			product.Desc3 = ""
		} else {
			product.Desc3 = "امتیاز " + messageCreator.CleaningString(desc3)
		}
	})

	err := c.Visit(url)
	if err != nil {
		println("visit error")
		return product, errors.New("visit error")
	}

	return product, nil

}
