package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"search/elastic"
	grpc2 "search/grpc"
)

type server struct {
	grpc2.UnimplementedSearchServer
}

func (s *server) Search(ctx context.Context, in *grpc2.SearchRequest) (*grpc2.SearchResponse, error) {
	log.Println(in.GetTerm())
	var products []*grpc2.Product
	for _, product := range elastic.SearchProduct(ctx, in.Term) {
		products = append(products, product.ToProto())
	}
	return &grpc2.SearchResponse{Products: products}, nil
}

func StartServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpc2.RegisterSearchServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
