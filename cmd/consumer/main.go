package main

import (
	"log"
	"os"
	"os/signal"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Shopify/sarama"
)

func main() {
	msgCount := 0
	brokerList := []string{"192.168.86.250:9092"}
	topic := "example"
	partition := int32(0)
	offset := sarama.OffsetOldest

	kingpin.Parse()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(topic, partition, offset)
	if err != nil {
		log.Panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				log.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				log.Println("Interrupt signal received... shuting down consumer")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Total message Processed - ", msgCount)
}
