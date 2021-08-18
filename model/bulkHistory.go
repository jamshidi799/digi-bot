package model

import "time"

type BulkHistory struct {
	ID              int `gorm:"primaryKey"`
	SourceProductID int
	CategoryID      int
	BrandID         int
	Price           int
	Date            time.Time
}
