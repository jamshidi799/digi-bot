package db

import (
	"digi-bot/model"
)

func GetAllProduct() []model.Product {
	var products []model.Product
	database.Find(&products)
	return products
}

func GetAllProductByUserId(userId int) []string {
	var products []string

	database.
		Select("product.name").
		Model(&model.Pivot{}).
		Joins("JOIN products product on product.id = pivots.product_id").
		//Preload("Pivots").
		Where("pivots.user_id = ?", userId).
		Find(&products)

	return products
}

func GetProductById(id int) model.Product {
	var product model.Product
	database.First(&product, id)
	return product
}

func SaveBulkHistories(histories *[]model.BulkHistory) {
	database.Create(histories)
}
