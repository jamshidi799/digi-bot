package model

import "github.com/jinzhu/gorm"

type PivotModel struct {
	gorm.Model
	ID        int
	UserId    int
	ProductId int
}
