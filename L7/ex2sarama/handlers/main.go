/// http:
/// 	GET on /rate?rate=value puts value to topic rate
/// 	GET on /total returns value from topic total
///
/// processes: listen for rate topic for value and puts the value to topic total

package main

import (
	"context"
	"fmt"
	"handlers/consumer"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

// KAFKA_BROKERS=localhost:9092 KAFKA_VERSION=4.0.0 go run main.go
// KAFKA_BROKERS=127.0.0.1:9092 go run main.go
// KAFKA_BROKERS=kafka:9092 KAFKA_VERSION=4.0.0 HTTP_PORT=8080 ./main

var (
	producer sarama.AsyncProducer
	totalCH  <-chan byte
	sigterm  = make(chan os.Signal, 1)

	VERSION = "0.0.0"
	// KAFKA_BROKERS = "localhost:9092"
	// KAFKA_VERSION = "4.0.0"
)

const (
	TOPIC_RATE  = "rate"
	TOPIC_TOTAL = "total"
	GROUP       = "totalGroup"
)

func main() {

	KAFKA_BROKERS := os.Getenv("KAFKA_BROKERS")
	KAFKA_VERSION := os.Getenv("KAFKA_VERSION")
	// HTTP_PORT := os.Getenv("HTTP_PORT")
	HTTP_PORT := "8080"

	log.Printf("--- start --- KAFKA_BROKERS(%v) KAFKA_VERSION(%v)\n", KAFKA_BROKERS, KAFKA_VERSION)

	//http
	mux := http.NewServeMux()
	mux.HandleFunc("/rate", PostRateHandler)
	mux.HandleFunc("/total", GetTotalHandler)
	mux.HandleFunc("/hc", HealthCheck)

	s := &http.Server{
		// Addr:    ":8080",
		Addr:    fmt.Sprintf(":%s", HTTP_PORT),
		Handler: mux,
	}

	//kafka stuff
	brokerList := []string{KAFKA_BROKERS}
	kafkaVersion, err := sarama.ParseKafkaVersion(KAFKA_VERSION)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	var reader *consumer.Reader
	reader, totalCH = consumer.NewReader([]string{TOPIC_TOTAL}, KAFKA_VERSION, brokerList, GROUP, &wg)
	reader.Start()

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
	// sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigterm
		log.Println("terminating: via signal")

		err := s.Shutdown(context.Background())
		if err != nil {
			log.Println("error on http shutdown: ", err)
		}
		log.Println("http server stopped")

		err = producer.Close()
		if err != nil {
			log.Println("error on producer close: ", err)
		}
		log.Println("producer closed")

		reader.Stop()
		log.Println("reader stopped")
	}()

	log.Println("http server started on ", HTTP_PORT)

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		// panic(fmt.Errorf("error on listen and serve: %v", err))
		log.Println(fmt.Errorf("error on http listen and serve: %v", err))
	}
	log.Println("http server stopped")

	wg.Wait()
	log.Println("main stopped")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {

}

func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetTotalHandler: got request")
	var (
		sum       int
		iter      int
		iter_prev int
	)

	stop := time.After(time.Millisecond * 500)

stop_label:
	for {
		select {
		case <-stop:
			log.Println("GetTotalHandler: stop")
			if iter != iter_prev {
				iter_prev = iter
				stop = time.After(time.Millisecond * 500)
			} else {
				break stop_label
			}
		case value, ok := <-totalCH:
			if !ok {
				log.Println("GetTotalHandler: no value in channel")
				break stop_label
			}
			sum += int(value)
			iter++
			// log.Println("GetTotalHandler: got value: ", value)
			log.Println("GetTotalHandler: got value, iter: ", iter)
		case <-r.Context().Done():
			log.Printf("GetTotalHandler: got context done err(%v)\n", r.Context().Err())
			break stop_label
		case <-sigterm:
			log.Printf("GetTotalHandler: got sigterm\n")
			break stop_label
		default:
			time.Sleep(time.Millisecond * 1)
		}
	}

	if iter != 0 {
		res := fmt.Sprintf("avg(%.2f); version(%v); iter(%v)) \n", float32(sum)/float32(iter), VERSION, iter)
		w.Write([]byte(res))
	} else {
		res := fmt.Sprintf("no data; version(%v)\n", VERSION)
		w.Write([]byte(res))
	}

	// value, ok := <-totalCH
	// if !ok {
	// 	log.Println("GetTotalHandler: no value in channel")
	// 	if sum != 0 {
	// 		res := float32(sum) / float32(iter)
	// 		w.Write([]byte(fmt.Sprintf("%.2f", res)))
	// 	} else {
	// 		w.Write([]byte("no data"))
	// 	}
	// 	return
	// }
	// log.Println("GetTotalHandler: got value: ", value)
	// sum += int(value)
	// iter++

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
