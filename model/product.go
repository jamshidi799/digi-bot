package model

import (
	"github.com/jinzhu/gorm"
)

type ProductModel struct {
	gorm.Model
	ID            int
	UserId        int
	Url           string
	Name          string
	AlternateName string
	Description   string
	Price         int
	OldPrice      int
	Image         string
}

func (productModel ProductModel) ToProduct() Product {
	return Product{
		Name:          productModel.Name,
		AlternateName: productModel.AlternateName,
		Url:           productModel.Url,
		Description:   productModel.Description,
		Price:         productModel.Price,
		OldPrice:      productModel.OldPrice,
		Images:        []string{productModel.Image},
	}
}

type Product struct {
	Name          string `json:"name"`
	AlternateName string `json:"alternateName"`
	Description   string `json:"description"`
	Url           string
	Offer         offer `json:"offers"`
	Price         int
	OldPrice      int
	Images        []string `json:"image"`
}

type offer struct {
	Price    int `json:"lowPrice"`
	OldPrice int `json:"highPrice"`
}

func (product Product) ToProductModel(userId int) ProductModel {
	return ProductModel{
		UserId:        userId,
		Url:           product.Url,
		Name:          product.Name,
		AlternateName: product.AlternateName,
		Description:   product.Description,
		Price:         product.Price,
		OldPrice:      product.OldPrice,
		Image:         product.Images[0],
	}
}
