package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Shopify/sarama"
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

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	admin, err := sarama.NewClusterAdmin([]string{broker}, config)
	if err != nil {
		log.Panic(err)
	}
	defer admin.Close()

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: 1,
	}, true)
	if err != nil && err != sarama.ErrTopicAlreadyExists {
		log.Panic(err)
	}

	producer, err := sarama.NewSyncProducer([]string{broker}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer producer.Close()

	for i := 0; i < messagesCount; i++ {
		msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder("Hello Kafka!!!")}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("FAILED to send message: %s\n", err)
		} else {
			log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
		}
	}

}
