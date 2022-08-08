package search

import (
	"context"
	"digi-bot/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func Query(term string) *SearchResponse {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewSearchClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := c.Search(ctx, &SearchRequest{Term: term})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return r
}

func (p *Product) ToDto() model.ProductDto {
	return model.ProductDto{
		Id:    int(p.Id),
		Name:  p.Title,
		Url:   p.Url,
		Price: int(p.Price),
		Sku:   p.Sku,
	}
}
