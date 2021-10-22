package model

import (
	"time"
)

type History struct {
	ID        int `gorm:"primaryKey"`
	ProductID int `gorm:"index"`
	Price     int
	Date      time.Time
}

type GraphData struct {
	Price int
	Date  time.Time
}
