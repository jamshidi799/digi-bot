package model

type Product struct {
	Id       int
	Name     string `json:"name"`
	Url      string
	Price    int
	OldPrice int
	Status   int
	Sku      string
}
