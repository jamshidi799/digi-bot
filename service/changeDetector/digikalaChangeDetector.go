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

// todo: just fetch digikala products
func (DigikalaChangeDetector) getProducts() []model.Product {
	return db.GetAllProduct()
}

func (dk *DigikalaChangeDetector) Detect(handler func(new model.ProductDto, old model.Product)) {
	for _, product := range dk.getProducts() {
		newProduct, err := dk.Crawl(product.Url)

		if err != nil {
			log.Println(err, newProduct)
			continue
		}

		handler(newProduct, product)

		time.Sleep(time.Second * 1)
	}
}
