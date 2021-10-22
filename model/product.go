package model

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	ID        int
	RealID    int
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
	Name     string `json:"name"`
	Url      string
	Price    int
	OldPrice int
	Status   int
	Desc1    string
	Desc2    string
	Desc3    string
}

func (product ProductDto) ToProduct() Product {
	return Product{
		Url:    product.Url,
		Name:   product.Name,
		Price:  product.Price,
		Status: product.Status,
	}
}
