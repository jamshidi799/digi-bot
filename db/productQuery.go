package db

import (
	"digi-bot/model"
	"log"
)

func GetAllProduct() []model.Product {
	var products []model.Product
	DB.Find(&products)
	return products
}

func GetAllProductByUserId(userId int) []string {
	var products []string

	DB.
		Select("product.name").
		Model(&model.Pivot{}).
		Joins("JOIN products product on product.id = pivots.product_id").
		Where("pivots.user_id = ?", userId).
		Find(&products)

	log.Println()

	return products
}

func GetProductById(id int) model.Product {
	var product model.Product
	DB.First(&product, id)
	return product
}
