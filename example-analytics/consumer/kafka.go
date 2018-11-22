package consumer

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func RunKafkaConsumer(consumerId string, broker string, group string, topic string,
	logger func(format string, v ...interface{}), process messageProcessator,
	quit chan struct{}) {

	config := sarama.NewConfig()
	config.ClientID = consumerId
	config.Consumer.Offsets.CommitInterval = time.Hour * 999999 //never
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		panic(err)
	}

	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalln(err)
	}

	toClose := make([]chan struct{}, 0)

	for _, partition := range partitions {
		go func(id int32) {
			myQuit := make(chan struct{})
			toClose = append(toClose, myQuit)
			partitionConsumer, err := consumer.ConsumePartition(topic, id, sarama.OffsetOldest)
			defer partitionConsumer.Close()

			if err != nil {
				panic(err)
			}
		ConsumerLoop:
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					//log.Printf("Consumed message offset %d\n", msg.Offset)
					process(msg.Value)
					//consumed++
				case <-myQuit:
					break ConsumerLoop
				}
			}

		}(partition)
	}

	<-quit
	logger("quit signal received in kafka-consumer, closing consumers")
	for _, toc := range toClose {
		toc <- struct{}{}
	}
	logger("done closing consumers")
}
