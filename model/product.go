package model

import (
	"github.com/jinzhu/gorm"
)

type ProductModel struct {
	gorm.Model
	ID    int
	Url   string
	Name  string
	Price int
}

func (productModel ProductModel) ToProduct() Product {
	return Product{
		Name:  productModel.Name,
		Url:   productModel.Url,
		Price: productModel.Price,
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

func (product Product) ToProductModel() ProductModel {
	return ProductModel{
		Url:   product.Url,
		Name:  product.Name,
		Price: product.Price,
	}
}
