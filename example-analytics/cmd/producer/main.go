package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bgadrian/kafka-playground/example-analytics/producer"

	"github.com/Shopify/sarama"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic> <partitionCount> <onlineUsersPeak> <eventsPerMinPerUser> \n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]
	partitions, err := strconv.Atoi(os.Args[3])

	if err != nil {
		log.Panic(err)
	}

	peak, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Panic(err)
	}

	eventsPerMinPerUser, err := strconv.Atoi(os.Args[5])
	if err != nil {
		log.Panic(err)
	}

	collector := make(chan producer.Event, peak)
	disp := buildDispatcher(peak, collector, eventsPerMinPerUser)
	kafkaProducer := buildKafka(broker, topic, partitions)

	//Tick of the TrafficPeak is 1 minute (simulating 1 hour/day with 1 minute)
	go func() {
		for ev := range collector {
			props := ev.Properties()
			propsBytes, _ := json.Marshal(props)
			log.Printf("user: %s, event: %s\n", props[producer.KeyUserID], ev.Name()) //, string(propsBytes)

			msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(propsBytes)}
			_, _, err := kafkaProducer.SendMessage(msg)
			if err != nil {
				log.Printf("FAILED to send message: %s\n", err)
			}
			//else {
			//	log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
			//}
		}
	}()
	disp.Start(time.Second * 5)

	term := make(chan os.Signal, 3)
	signal.Notify(term, syscall.SIGTERM)
	<-term //wait for ctrl+c

	disp.Stop()
	close(collector)
	kafkaProducer.Close()
}

func buildKafka(broker, topic string, partitions int) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer([]string{broker}, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return producer
}

func buildDispatcher(peak int, collector chan producer.Event, eventsPerMinPerUser int) *producer.Dispatcher {
	//TODO move this into a YAML
	eventTemplates := []producer.Event{
		//if you modify these update the README too
		//for the supported {variable} see https://github.com/bgadrian/fastfaker/blob/master/TEMPLATE_VARIABLES.md
		producer.NewSimpleEvent("pageView", map[string]string{"page": "/{bs}/#/"}),
		producer.NewSimpleEvent("click", map[string]string{"element": "{hackerabbreviation}"}),
		producer.NewSimpleEvent("error", map[string]string{"location": "{hackeringverb}:{hackernoun}:###"}),
		producer.NewSimpleEvent("addToCart", map[string]string{
			"name":   "{carmaker}",
			"color":  " {safecolor}",
			"price":  "{uint8}.#",
			"amount": "#",
		}),
		producer.NewSimpleEvent("buy", map[string]string{
			"name":   "{carmaker}",
			"color":  " {safecolor}",
			"price":  "{uint8}.#",
			"amount": "#",
		}),
	}

	gen := producer.NewRandomEventGenerator(eventTemplates)
	traff := producer.NewTrafficNormalDay(peak)

	//a pool of user IDs to be reused, with future sessions
	//simulating reocurring users
	oldUsersRevisiting := make(chan string, peak/10)

	queen := func() producer.User {
		//each user will generate a batch of events each random( minMs and maxMs)
		second := 60000
		minMs := 1 * second
		maxMs := 3 * second
		user := producer.NewSimpleUser(eventsPerMinPerUser, gen, collector, minMs, maxMs)
		userID, _ := user.Property(producer.KeyUserID)

		//40% chance to be an old user, if we have any
		if len(oldUsersRevisiting) > 0 && rand.Intn(10) < 4 {
			user.SetProperty(producer.KeyUserID, <-oldUsersRevisiting)
		} else if rand.Intn(10) < 1 { //10% anonymous users
			//10% anonymous users
			user.SetProperty(producer.KeyUserID, "")
		}

		//add the user to the OLD pool
		if len(oldUsersRevisiting) < cap(oldUsersRevisiting) {
			oldUsersRevisiting <- userID
		}

		log.Printf("new user: %s\n", userID)
		return user
	}

	return producer.NewDispatcher(traff, queen)
}
