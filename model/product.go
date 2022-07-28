package model

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	ID        int
	Sku       string
	Url       string
	Name      string
	Price     int
	Status    int
	Pivots    []Pivot
	Histories []History `gorm:"foreignKey:ProductID"`
}

func (product Product) ToDto() ProductDto {
	return ProductDto{
		Name:   product.Name,
		Url:    product.Url,
		Price:  product.Price,
		Status: product.Status,
	}
}

type ProductDto struct {
	Id       int
	Name     string `json:"name"`
	Url      string
	Price    int
	OldPrice int
	Status   int
	Sku      string
}

func (product ProductDto) ToProduct() Product {
	return Product{
		Url:    product.Url,
		Name:   product.Name,
		Price:  product.Price,
		Status: product.Status,
		Sku:    product.Sku,
	}
}
