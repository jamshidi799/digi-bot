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
	Image     string
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
		Image:  product.Image,
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
	Image    string
	Discount int
}

func (product ProductDto) ToProduct() Product {
	return Product{
		Url:    product.Url,
		Name:   product.Name,
		Price:  product.Price,
		Status: product.Status,
		Sku:    product.Sku,
		Image:  product.Image,
	}
}
