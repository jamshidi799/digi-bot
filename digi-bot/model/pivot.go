package model

import "github.com/jinzhu/gorm"

type Pivot struct {
	gorm.Model
	ID        int
	UserID    int
	ProductID int
	Discount  int
}

type UserIdAndDiscountDto struct {
	UserId   int
	Discount int
}
