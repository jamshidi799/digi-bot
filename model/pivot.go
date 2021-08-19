package model

import "github.com/jinzhu/gorm"

type Pivot struct {
	gorm.Model
	ID                  int
	UserID              int
	ProductID           int
	NotificationSetting int8
}
