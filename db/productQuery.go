package db

import (
	"digi-bot/model"
)

func GetAllProduct() []model.ProductModel {
	var products []model.ProductModel
	DB.Find(&products)
	return products
}

func GetAllProductByUserId(userId int) []string {
	var products []string

	DB.
		Select("product.name").
		Model(&model.PivotModel{}).
		Joins("JOIN product_models product on product.id = pivot_models.product_id").
		Where("pivot_models.user_id = ?", userId).
		Find(&products)

	return products
}

func GetProductById(id int) model.ProductModel {
	var product model.ProductModel
	DB.First(&product, id)
	return product
}
