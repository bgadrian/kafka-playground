package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bgadrian/kafka-playground/example-analytics/consumer"
)

var logger = log.Printf

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic> \n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]
	group := "dashboard"
	consumerId := group + "_main"

	closeKafka := make(chan struct{})

	go consumer.RunKafkaConsumer(consumerId, broker, group, topic, logger, consumer.ProcessMessage, closeKafka)

	//blocks until a key is pressed
	consumer.Run()

	//blocks until kafka is closed
	close(closeKafka)
	os.Exit(0)
}
