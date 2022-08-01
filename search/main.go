package main

import (
	"search/consumer"
	"search/elastic"
)

func main() {
	elastic.InitElasticsearch()
	consumer.Consume()
}
