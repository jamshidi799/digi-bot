package entity

import (
	"digi-bot/model"
	"time"
)

type BatchCrawlResult struct {
	Page       int       `njson:"data.trackerData.page"`
	Pages      int       `njson:"data.trackerData.pages"`
	CategoryId int       `njson:"data.trackerData.categoryId"`
	FoundItems int       `njson:"data.trackerData.foundItems"`
	Products   []Product `njson:"data.trackerData.products"`
}

type Product struct {
	ProductId  int `njson:"product_id"`
	CategoryId int `njson:"category_id"`
	BrandId    int `njson:"brand_id"`
	Price      int `njson:"selling_price"`
}

func (product Product) ToBulkHistoryModel() model.BulkHistory {
	return model.BulkHistory{
		SourceProductID: product.ProductId,
		CategoryID:      product.CategoryId,
		BrandID:         product.BrandId,
		Price:           product.Price,
		Date:            time.Now(),
	}
}

func (data BatchCrawlResult) ToBatchBulkHistory() []model.BulkHistory {
	var bulk []model.BulkHistory
	for _, product := range data.Products {
		//if product.Price == 0 {
		//	return bulk
		//}
		bulk = append(bulk, product.ToBulkHistoryModel())
	}
	return bulk
}
