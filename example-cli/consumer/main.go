package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <count> <topics..> \n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	group := os.Args[2]
	consumerCount, err := strconv.Atoi(os.Args[3])
	topics := os.Args[4:]

	if err != nil {
		log.Panic(err)
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	consumersQuitChannels := make([]chan bool, 0, consumerCount)
	wg := sync.WaitGroup{}
	wg.Add(consumerCount)

	log.Printf("starting %d consumers\n", consumerCount)

	for i := 0; i < consumerCount; i++ {
		qchForI := make(chan bool, 2)
		consumersQuitChannels = append(consumersQuitChannels, qchForI)

		go func(i int, quitChannel chan bool) {
			consumerId := fmt.Sprintf("%s_%d", group, i)

			c, err := kafka.NewConsumer(&kafka.ConfigMap{
				"client.id":                       fmt.Sprintf(consumerId, i),
				"bootstrap.servers":               broker,
				"group.id":                        group,
				"session.timeout.ms":              6000,
				"go.events.channel.enable":        true,
				"go.application.rebalance.enable": true,
				"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
				"fetch.wait.max.ms":               10,
				"fetch.error.backoff.ms":          50,
			})

			if err != nil {
				log.Panic(err)
			}

			err = c.SubscribeTopics(topics, nil)
			defer c.Close()

			fmt.Printf("connecting to %s, as %s for topics %v \n",
				broker, consumerId, topics)

			process(consumerId, qchForI, c)

			_, err = c.Commit() //save our work!
			if err != nil {
				log.Printf("err %s \n", err.Error())
			}
			log.Printf("%s finished \n", consumerId)
			wg.Done()
		}(i, qchForI)
	}

	//first we wait for the sigterm
	<-quitChannel
	log.Printf("closing %d consumers\n", consumerCount)

	for _, cq := range consumersQuitChannels {
		cq <- true
	}

	//second we wait for the consumers to finish processing
	wg.Wait()
	fmt.Println("the end")
}

func process(consumerId string, quitChannel chan bool, c *kafka.Consumer) {
	for {
		select {
		case <-quitChannel:
			return
		case ev := <-c.Events():
			switch e := ev.(type) {

			case *kafka.Message:
				fmt.Printf("%s Message on %s:\n%s\n",
					consumerId, e.TopicPartition, string(e.Value))

			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%s %v\n", consumerId, e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%s %v\n", consumerId, e)
				c.Unassign()
			case kafka.PartitionEOF:
				fmt.Printf("%s Reached %v\n", consumerId, e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%s Error: %v\n", consumerId, e)
			}
		}
	}
}
