package db

import "digi-bot/model"

func GetAllProduct() []model.ProductModel {
	var objects []model.ProductModel
	DB.Find(&objects)
	return objects
}
