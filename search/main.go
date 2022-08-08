package main

import (
	"search/consumer"
	"search/elastic"
)

func main() {
	elastic.InitElasticsearch()
	go consumer.Consume()

	StartServer()
}
