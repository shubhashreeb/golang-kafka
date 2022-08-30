package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func main() {

	brokerList := []string{"192.168.86.250:9092"}
	topic := "example"
	totalMessage := 100000

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()

	abc := func(index int, msg *sarama.ProducerMessage) {
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Message Seq - (%d) is stored in topic(%s)/partition(%d)/offset(%d)\n", index, topic, partition, offset)
	}

	for i := 0; i < totalMessage; i++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("Sample Message :: %d", i)),
		}
		go abc(i, msg)
	}
	time.Sleep(10 * time.Second)
}
