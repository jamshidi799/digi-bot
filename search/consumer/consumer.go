package consumer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Consumer interface {
	DoConsume(message *kafka.Message)
	Topic() string
}

var consumers = []Consumer{productConsumer{"products"}}

func Consume() {

	for _, consumer := range consumers {
		c := getConsumer(consumer.Topic())

		// Set up a channel for handling Ctrl-C, etc
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

		run := true
		for run == true {
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				run = false
			default:
				ev, err := c.ReadMessage(100 * time.Millisecond)
				if err != nil {
					// Errors are informational and automatically handled by the consumer
					continue
				}

				consumer.DoConsume(ev)
			}
		}

		// todo
		defer c.Close()
	}

}

func getConsumer(topic string) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "pkc-ymrq7.us-east-2.aws.confluent.cloud:9092",
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     "ADO7I5Q7MDGYYCT2",
		"sasl.password":     "1Qzl10sz4hYUqzuMxLLpNBJrCylEXe3Xz311D8PPjwthsnm0O/+oHszk50P8E04e",
		"group.id":          "test",
		"auto.offset.reset": "earliest",
	})

	log.Printf("kafka consumer connected to %s", topic)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	return c
}
