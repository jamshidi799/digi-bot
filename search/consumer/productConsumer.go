package consumer

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"search/elastic"
	"search/model"
)

type productConsumer struct {
	topic string
}

func (p productConsumer) Topic() string {
	return p.topic
}

func (p productConsumer) DoConsume(message *kafka.Message) {
	var product model.Product

	if err := json.Unmarshal(message.Value, &product); err != nil {
		log.Println(err)
		return
	}
	log.Printf("handling product with id: %d", product.Id)
	go elastic.AddProduct(product)
}
