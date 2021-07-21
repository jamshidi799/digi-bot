package entity

type Object struct {
	Name          string   `json:"name"`
	AlternateName string   `json:"alternateName"`
	Description   string   `json:"description"`
	Price         int      `json:"offers.lowPrice"`
	OldPrice      int      `json:"offers.highPrice"`
	Images        []string `json:"image"`
}
