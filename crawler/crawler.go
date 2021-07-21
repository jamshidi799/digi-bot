package crawler

import (
	"digi-bot/entity"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
)

func Crawl() {
	c := colly.NewCollector(
		colly.AllowedDomains("digikala.com", "www.digikala.com"),
	)

	c.OnHTML("script[type=\"application/ld+json\"]", func(e *colly.HTMLElement) {
		var obj entity.Object
		err := json.Unmarshal([]byte(e.Text), &obj)
		if err != nil {
			println("err")
		}

		if obj.Name != "" {
			fmt.Printf("%+v\n", obj)
		}
	})

	url := "https://www.digikala.com/product/dkp-4892730/%D9%84%D9%BE-%D8%AA%D8%A7%D9%BE-14-%D8%A7%DB%8C%D9%86%DA%86%DB%8C-%D9%87%D9%88%D8%A2%D9%88%DB%8C-%D9%85%D8%AF%D9%84-matebook-d14-a"
	err := c.Visit(url)
	if err != nil {
		println("visit error")
	}

}
