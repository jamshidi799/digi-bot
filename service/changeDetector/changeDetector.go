package changeDetector

import (
	"digi-bot/model"
	"digi-bot/service/crawler"
)

type ChangeDetector interface {
	crawler.Crawler
	Detect(handler func(new model.ProductDto, old model.Product))
	getProducts() []model.Product
}
