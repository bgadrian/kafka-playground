package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic> <partitions> <msg count>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]
	partitions, err := strconv.Atoi(os.Args[3])
	messagesCount, _ := strconv.Atoi(os.Args[4])

	if err != nil {
		log.Panic(err)
	}

	//https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	//we try to "kill" the kafkas performance by sending one message
	// at a time, for demo purposes
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":       broker,
		"queue.buffering.max.ms":  5,
		"go.events.channel.size":  100,
		"go.produce.channel.size": 10,
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	createTopic(broker, topic, partitions, 1)
	fmt.Printf("Created Producer %v\n", p)

	wg := &sync.WaitGroup{}
	wg.Add(messagesCount)

	go report(p, wg)
	write(p, "Hello kafka!", topic, messagesCount)

	wg.Wait()
	p.Close()
}

var write = func(p *kafka.Producer, msg, topic string, count int) {

	for i := 0; i < count; i++ {
		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(msg),
			Key:   []byte(time.Now().UTC().String()),
		}

		//blocks when is full, see go.produce.channel.size
		p.ProduceChannel() <- message
	}

}

var report = func(p *kafka.Producer, wg *sync.WaitGroup) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		default:
			fmt.Printf("Ignored event: %s\n", ev)
		}
		wg.Done()
	}
}
