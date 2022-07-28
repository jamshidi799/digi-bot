package changeDetector

import (
	"digi-bot/model"
	"digi-bot/model/db"
	"digi-bot/service/crawler"
	"log"
	"time"
)

type TorobChangeDetector struct {
	crawler.TorobCrawler
}

func (detector *TorobChangeDetector) getProducts() []model.Product {
	return db.GetAllProductByDomain(detector.GetDomain())
}

func (detector *TorobChangeDetector) Detect(handler func(new model.ProductDto, old model.Product)) {
	for _, product := range detector.getProducts() {
		newProduct, err := detector.Crawl(product.Url)

		if err != nil {
			log.Println(err, newProduct)
			continue
		}

		handler(newProduct, product)

		time.Sleep(time.Second * 1)
	}
}
