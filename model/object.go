package model

import (
	"github.com/jinzhu/gorm"
)

type ObjectModel struct {
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

type Object struct {
	Name          string   `json:"name"`
	AlternateName string   `json:"alternateName"`
	Description   string   `json:"description"`
	Price         int      `json:"offers.lowPrice"`
	OldPrice      int      `json:"offers.highPrice"`
	Images        []string `json:"image"`
}

func (obj Object) ToObjectModel(userId int, url string) ObjectModel {
	return ObjectModel{
		UserId:        userId,
		Url:           url,
		Name:          obj.Name,
		AlternateName: obj.AlternateName,
		Description:   obj.Description,
		Price:         obj.Price,
		OldPrice:      obj.OldPrice,
		Image:         obj.Images[0],
	}
}
