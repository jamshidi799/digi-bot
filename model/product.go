package model

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	ID        int
	Url       string
	Name      string
	Price     int
	Pivots    []Pivot
	Histories []History `gorm:"foreignKey:ProductID"`
}

func (product Product) ToDto() ProductDto {
	return ProductDto{
		Name:  product.Name,
		Url:   product.Url,
		Price: product.Price,
	}
}

type ProductDto struct {
	Name     string `json:"name"`
	Url      string
	Price    int
	OldPrice int
	Desc1    string
	Desc2    string
	Desc3    string
}

func (product ProductDto) ToProduct() Product {
	return Product{
		Url:   product.Url,
		Name:  product.Name,
		Price: product.Price,
	}
}
