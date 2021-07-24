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

func (objModel ObjectModel) ToObject() Object {
	return Object{
		Name:          objModel.Name,
		AlternateName: objModel.AlternateName,
		Url:           objModel.Url,
		Description:   objModel.Description,
		Price:         objModel.Price,
		OldPrice:      objModel.OldPrice,
		Images:        []string{objModel.Image},
	}
}

type Object struct {
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
