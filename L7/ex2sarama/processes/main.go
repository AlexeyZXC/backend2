package main

import (
	"context"
	"log"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/kafkapubsub"
)

var (
	ratesTopic *pubsub.Topic
	totalTopic *pubsub.Topic

	ratesSub *pubsub.Subscription
	totalSub *pubsub.Subscription
)

const (
	// "kafka://process?topic=rates")
	// ratesAddr = "kafka://rates"
	ratesAddr = "kafka://rates?topic=rates"
	totalAddr = "kafka://total"
)

// KAFKA_BROKERS=kafka:9092 go run main.go
// KAFKA_BROKERS=localhost:9092 go run main.go

// var (
// 	connFunc = func(network, addr string) (radix.Conn, error) {
// 		return radix.Dial(network, addr,
// 			radix.DialTimeout(timeout),
// 		)
// 	}
// )

// func subscription() (*pubsub.Subscription, error) {
// 	if sub != nil {
// 		return sub, nil
// 	}
// 	var err error
// 	sub, err = pubsub.OpenSubscription(context.Background(),
// 		"kafka://process?topic=rates")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sub, nil
// }

func main() {
	for {
		s, err := getSubscription(ratesAddr)
		if err != nil {
			log.Fatal("getSubscription: %v", err)
		}

		msg, err := s.Receive(context.Background())
		if err != nil {
			log.Fatal("Receive: %v", err)
		}

		// err = storage().Do(radix.Cmd(nil, "LPUSH", "result", string(msg.Body)))
		// if err != nil {
		// 	log.Println(err)
		// }
		// if rand.Float64() < .05 {
		// 	_ = storage().Do(radix.Cmd(nil, "LTRIM", "result", "0", "9"))
		// }
		log.Println("got body: ", string(msg.Body))

		t := getTopic(totalAddr)

		err = t.Send(context.Background(), &pubsub.Message{
			Body: msg.Body,
		})
		if err != nil {
			log.Fatal("Send: %v", err)
		}
		log.Println("sent message: ", string(msg.Body))

		msg.Ack()
	}
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
		log.Fatal("getSubscription: fail to OpenSubscription: ", err)
		// return nil, err
	}
	return sub, nil
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
		panic(err)
	}
	return t
}
