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
		var result map[string]interface{}
		err := json.Unmarshal([]byte(e.Text), &result)
		if err != nil {
			println("err")
		}
		obj := entity.Object{}
		if val, ok := result["name"].(string); ok {
			obj.Name = val
		}
		if val, ok := result["alternateName"].(string); ok {
			obj.AlternateName = val
		}
		if val, ok := result["description"].(string); ok {
			obj.Description = val
		}
		if val, ok := result["offers"].(map[string]interface{}); ok {
			obj.Price = int(val["lowPrice"].(float64))
		}
		if val, ok := result["offers"].(map[string]interface{}); ok {
			obj.OldPrice = int(val["highPrice"].(float64))
		}

		obj.Images = []string{}
		if val, ok := result["image"].([]interface{}); ok {
			for _, value := range val {
				obj.Images = append(obj.Images, value.(string))
			}
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
