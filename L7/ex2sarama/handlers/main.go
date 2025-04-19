package main

import (
	"fmt"
	"handlers/consumer"
	"log"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
)

// KAFKA_BROKERS=localhost:9092 go run main.go
// KAFKA_BROKERS=127.0.0.1:9092 go run main.go

var (
	producer sarama.AsyncProducer
)

const (
	KAFKA_BROKERS = "localhost:9092"
	KAFKA_VERSION = "4.0.0"
	TOPIC_RATE    = "rate"
	GROUP         = "group1"
)

func init() {
	// router.
	// 	HandleFunc("/rate", PostRateHandler).
	// 	Methods(http.MethodPost)
	// router.
	// 	HandleFunc("/total", GetTotalHandler).
	// 	Methods(http.MethodGet)

	http.HandleFunc("/rate", PostRateHandler)
	http.HandleFunc("/total", GetTotalHandler)
}

func main() {

	brokerList := []string{KAFKA_BROKERS}
	kafkaVersion, err := sarama.ParseKafkaVersion(KAFKA_VERSION)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	reader, _ := consumer.NewReader([]string{TOPIC_RATE}, KAFKA_VERSION, brokerList, GROUP)
	go reader.Start()

	config := sarama.NewConfig()
	config.Version = kafkaVersion
	// tlsConfig := createTlsConfiguration()
	// if tlsConfig != nil {
	// 	config.Net.TLS.Enable = true
	// 	config.Net.TLS.Config = tlsConfig
	// }
	config.Producer.RequiredAcks = sarama.WaitForLocal // Only wait for the leader to ack
	// config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	// config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err = sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	log.Println("http server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}

func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
	// var rates []string

	// sub, err := getSubscription(totalAddr)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	// defer cancel()

	// msg, err := sub.Receive(ctx)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// log.Printf("GetTotalHandler: msgBody(%v) ", string(msg.Body))

	// var sum int
	// for _, rate := range rates {
	// 	v, err := strconv.Atoi(rate)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	sum += v
	// }
	// result := float64(sum) / float64(len(rates))
	// _, _ = w.Write([]byte(fmt.Sprintf("%.2f", result)))
}

func PostRateHandler(w http.ResponseWriter, r *http.Request) {
	rateStr := r.FormValue("rate")
	log.Printf("PostRateHandler: got request: rate: %v\n", rateStr)
	var (
		err  error
		rate int
	)

	if rate, err = strconv.Atoi(rateStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	producer.Input() <- &sarama.ProducerMessage{
		Topic: TOPIC_RATE,
		Value: sarama.ByteEncoder{byte(rate)},
	}

	log.Printf("PostRateHandler: send ok: %v\n", rateStr)
}
