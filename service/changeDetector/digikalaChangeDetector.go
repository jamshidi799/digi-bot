package changeDetector

import (
	"digi-bot/model"
	"digi-bot/model/db"
	"digi-bot/service/crawler"
	"log"
	"time"
)

type DigikalaChangeDetector struct {
	crawler.DigikalaCrawler
}

func (detector *DigikalaChangeDetector) getProducts() []model.Product {
	return db.GetAllProductByDomain(detector.GetDomain())
}

func (detector *DigikalaChangeDetector) Detect(handler func(new model.ProductDto, old model.Product)) {
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
