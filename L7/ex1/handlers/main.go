package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gocloud.dev/pubsub"

	_ "gocloud.dev/pubsub/kafkapubsub"
)

// KAFKA_BROKERS=localhost:9092 go run main.go
// KAFKA_BROKERS=127.0.0.1:9092 go run main.go

var (
	ratesTopic *pubsub.Topic
	totalTopic *pubsub.Topic

	ratesSub *pubsub.Subscription
	totalSub *pubsub.Subscription
)

const (
	ratesAddr = "kafka://rates"
	totalAddr = "kafka://total"
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
	// if err := web.ListenAndServe(); err != http.ErrServerClosed {
	// 	panic(fmt.Errorf("error on listen and serve: %v", err))
	// }

	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}

// func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
// 	var rates []string
// 	err := storage().Do(radix.Cmd(&rates, "LRANGE", "result", "0", "10"))
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	if len(rates) == 0 {
// 		w.WriteHeader(http.StatusServiceUnavailable)
// 		return
// 	}
// 	var sum int
// 	for _, rate := range rates {
// 		v, err := strconv.Atoi(rate)
// 		if err != nil {
// 			continue
// 		}
// 		sum += v
// 	}
// 	result := float64(sum) / float64(len(rates))
// 	_, _ = w.Write([]byte(fmt.Sprintf("%.2f", result)))
// }

func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
	// var rates []string

	// err := storage().Do(radix.Cmd(&rates, "LRANGE", "result", "0", "10"))
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// if len(rates) == 0 {
	// 	w.WriteHeader(http.StatusServiceUnavailable)
	// 	return
	// }

	sub, err := getSubscription(totalAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	msg, err := sub.Receive(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("GetTotalHandler: msgBody(%v) ", string(msg.Body))

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
	rate := r.FormValue("rate")
	log.Printf("PostRateHandler: got request: rate: %v\n", rate)
	if _, err := strconv.Atoi(rate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := getTopic(ratesAddr).Send(context.Background(), &pubsub.Message{
		Body: []byte(rate),
	})
	if err != nil {
		log.Println("getTopic err: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("PostRateHandler: send ok: %v", rate)
}

func getTopic(topicAddr string) *pubsub.Topic {
	var err error
	var t *pubsub.Topic

	switch topicAddr {
	case ratesAddr:
		t = ratesTopic
	case totalAddr:
		t = totalTopic
	default:
		log.Fatal("getTopic: wrong topic addr: %v", topicAddr)
	}

	if t != nil {
		return t
	}

	t, err = pubsub.OpenTopic(context.Background(), topicAddr)
	if err != nil {
		log.Fatal("getTopic: fail to OpenTopic: ", err)
	}
	return t
}

func getSubscription(addr string) (*pubsub.Subscription, error) {
	var (
		sub *pubsub.Subscription
		err error
	)

	switch addr {
	case ratesAddr:
		sub = ratesSub
	case totalAddr:
		sub = totalSub
	default:
		log.Fatal("getSubscription: wrong sub addr: %v", addr)
	}

	if sub != nil {
		return sub, nil
	}

	sub, err = pubsub.OpenSubscription(context.Background(),
		// "kafka://process?topic=rates")
		addr)
	if err != nil {
		log.Fatal("getSubscription: fail to OpenSubscription: %v", err)
		// return nil, err
	}
	return sub, nil
}

// func storage() *radix.Pool {
// 	if s != nil {
// 		return s
// 	}
// 	var err error
// 	s, err = radix.NewPool("tcp", "redis:6379", 1, radix.PoolConnFunc(connFunc))
// 	if err != nil {
// 		panic(err)
// 	}
// 	return s
// }
