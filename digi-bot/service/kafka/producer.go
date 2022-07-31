package kafka

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
)

const (
	flushTimeout = 10000
)

var producer *kafka.Producer

func InitProducer() {
	log.Println("kafka producer try to connect")

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "pkc-ymrq7.us-east-2.aws.confluent.cloud:9092",
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     "ADO7I5Q7MDGYYCT2",
		"sasl.password":     "1Qzl10sz4hYUqzuMxLLpNBJrCylEXe3Xz311D8PPjwthsnm0O/+oHszk50P8E04e",
		"acks":              "all",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	log.Println("kafka producer connected")

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	producer = p
}

func Send(topic string, key string, data []byte) {
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          data,
	}, nil)

	if err != nil {
		log.Println(err)
	}
}

func FlushAndClose() {
	producer.Flush(flushTimeout)
	producer.Close()
}
