package model

import (
	"github.com/jinzhu/gorm"
)

type ProductModel struct {
	gorm.Model
	ID       int
	UserId   int
	Url      string
	Name     string
	Price    int
	OldPrice int
}

func (productModel ProductModel) ToProduct() Product {
	return Product{
		Name:     productModel.Name,
		Url:      productModel.Url,
		Price:    productModel.Price,
		OldPrice: productModel.OldPrice,
	}
}

type Product struct {
	Name     string `json:"name"`
	Url      string
	Price    int
	OldPrice int
	Desc1    string
	Desc2    string
	Desc3    string
}

func (product Product) ToProductModel(userId int) ProductModel {
	return ProductModel{
		UserId:   userId,
		Url:      product.Url,
		Name:     product.Name,
		Price:    product.Price,
		OldPrice: product.OldPrice,
	}
}
