package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type PriceData struct {
	PriceUsd string
	Time     uint64
}

type ApiResponse struct {
	Data      []PriceData
	Timestamp uint64
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka1:19092",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	topic := "bitcoin"

	// basically just logs message deliveries asynchronously
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key=%-10s, value=%s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			case *kafka.Error:
				fmt.Printf("Error: %s", ev)
			default:
				fmt.Printf("Kafka event: %s", ev)
			}

		}
	}() // sure

	// now here is where we need to fetch from API and send it
	url := "http://api.coincap.io/v2/assets/bitcoin/history?interval=m1"

	client := &http.Client{}

	res, err := client.Get(url)

	if err != nil {
		fmt.Printf("Failed to obtain response: %s", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s", err)
		os.Exit(1)
	}

	var apiResponse ApiResponse

	json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Printf("Failed to decode JSON: %s", err)
	}

	println("Read to send data...")
	for _, v := range apiResponse.Data {
		// send data
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(fmt.Sprint(v.Time)),
			Value:          []byte(v.PriceUsd),
		}, nil)
		if err != nil {
			fmt.Printf("Failed to produce message: %s", err)
			os.Exit(1)
		}
		time.Sleep(time.Second)
	}

	p.Close()
}
