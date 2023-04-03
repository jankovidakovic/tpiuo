package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func initializeProducer() *kafka.Producer {
	hostname, _ := os.Hostname()
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "bla bla",
		"client.id":         hostname,
		"acks":              "all",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	return p
}
