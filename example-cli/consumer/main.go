package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic> \n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	consumersQuitChannels := make([]chan struct{}, 0)

	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalln(err)
	}
	for _, partition := range partitions {
		go func(id int32) {
			myQuit := make(chan struct{})
			consumersQuitChannels = append(consumersQuitChannels, myQuit)
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			defer partitionConsumer.Close()

			if err != nil {
				panic(err)
			}
		ConsumerLoop:
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					log.Printf("Message read: %s", string(msg.Value))
				case <-myQuit:
					break ConsumerLoop
				}
			}

		}(partition)
	}

	//first we wait for the sigterm
	<-quitChannel
	log.Printf("closing consumers\n")

	for _, cq := range consumersQuitChannels {
		cq <- struct{}{}
	}

	//second we wait for the consumers to finish processing
	fmt.Println("the end")
}
