package model

import (
	"search/grpc"
)

type Product struct {
	Id       int
	Name     string
	Url      string
	Price    int
	OldPrice int
	Status   int
	Sku      string
	Image    string
}

func (p *Product) ToProto() *grpc.Product {
	return &grpc.Product{
		Id:    int32(p.Id),
		Title: p.Name,
		Price: int64(p.Price),
		Sku:   p.Sku,
		Image: p.Image,
		Url:   p.Url,
	}
}
