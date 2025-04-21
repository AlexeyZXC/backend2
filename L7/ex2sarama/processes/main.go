package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"

	"handlers/consumer"
)

var (
	producer sarama.AsyncProducer
	rateCH   <-chan byte
	totalCH  chan<- byte
)

const (
	KAFKA_BROKERS = "localhost:9092"
	KAFKA_VERSION = "4.0.0"
	TOPIC_RATE    = "rate"
	TOPIC_TOTAL   = "total"
	GROUP         = "rateGroup"
)

func main() {
	//kafka stuff
	brokerList := []string{KAFKA_BROKERS}
	kafkaVersion, err := sarama.ParseKafkaVersion(KAFKA_VERSION)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	var reader *consumer.Reader
	reader, rateCH = consumer.NewReader([]string{TOPIC_RATE}, KAFKA_VERSION, brokerList, GROUP, &wg)
	reader.Start()

	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Producer.RequiredAcks = sarama.WaitForLocal // Only wait for the leader to ack
	// config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	// config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err = sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	wg.Add(1)
	go func() {
		defer func() {
			log.Println("producer reading errors stopped")
			wg.Done()
		}()

		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	//stop stuff
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	//process stuff
	wg.Add(1)
	go func() {
		defer func() {
			log.Println("process stuff stopped")
			wg.Done()
		}()

	stop:
		for {
			select {
			case <-sigterm:
				log.Println("terminating: via signal")

				err = producer.Close()
				if err != nil {
					log.Println("error on producer close: ", err)
				}

				reader.Stop()
				break stop

			case b := <-rateCH:
				log.Printf("process got rate(%v)\n", b)
				producer.Input() <- &sarama.ProducerMessage{Topic: TOPIC_TOTAL, Value: sarama.ByteEncoder{b}}
			}
		}
	}()

	wg.Wait()
	log.Println("main processes done")
}
