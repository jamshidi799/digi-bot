package model

import "time"

type BulkHistory struct {
	ID              int `gorm:"primaryKey"`
	SourceProductID int `gorm:"index"`
	CategoryID      int
	BrandID         int
	Price           int
	Date            time.Time
}
