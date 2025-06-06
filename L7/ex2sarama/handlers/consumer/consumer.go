package consumer

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Reader struct {
	Topics []string
	// Msgs         <-chan byte
	consumer     *consumer
	KafkaVersion sarama.KafkaVersion
	cancelFunc   context.CancelFunc
	brokers      []string
	group        string
	wgInternal   sync.WaitGroup
	wgExternal   *sync.WaitGroup
}

func NewReader(topic []string, kafkaVersion string, brokers []string, group string, wg *sync.WaitGroup) (*Reader, <-chan byte) {
	version, err := sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	msgs := make(chan byte, 1000)

	return &Reader{
			Topics: topic,
			// Msgs:  msgs,
			consumer:     newConsumer(msgs),
			KafkaVersion: version,
			brokers:      brokers,
			group:        group,
			wgExternal:   wg,
		},
		msgs
}

func (r *Reader) Start() {
	config := sarama.NewConfig()
	config.Version = r.KafkaVersion

	var (
		ctx context.Context
	)

	ctx, r.cancelFunc = context.WithCancel(context.Background())

	client, err := sarama.NewConsumerGroup(r.brokers, r.group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	r.wgInternal.Add(1)
	go func() {
		defer func() {
			log.Println("consumer loop done")
			r.wgInternal.Done()
		}()

		for {
			log.Println("consumer loop")
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, r.Topics, r.consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					log.Println("sarama.ErrClosedConsumerGroup")
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println("consumer loop: context errror: ", ctx.Err().Error())
				return
			}
			// consumer.ready = make(chan bool)
			r.consumer.ready = make(chan bool)
		}
	}()

	<-r.consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")
}

func (r *Reader) Stop() {
	log.Println("consumer stopping...")
	r.cancelFunc()
	log.Println("consumer stopping...0")
	close(r.consumer.Msgs)
	log.Println("consumer stopping...1")
	r.wgInternal.Wait()
	log.Println("consumer stopping...2")
	r.wgExternal.Done()
	log.Println("consumer stopped")
}

// consumer represents a Sarama consumer group consumer
type consumer struct {
	ready  chan bool
	Msgs   chan byte
	doneCH chan struct{}
}

func newConsumer(msgs chan byte) *consumer {
	return &consumer{
		ready:  make(chan bool),
		Msgs:   msgs,
		doneCH: make(chan struct{}),
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("consumer setup")
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("consumer cleanup")
	// consumer.doneCH <- struct{}{}
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("ConsumeClaim: message channel was closed")
				return nil
			}
			log.Printf("ConsumeClaim: Message claimed: valueLen(%v) timestamp(%v) topic(%v) Partition(%v) Offset(%v) Headers(%+v)\n",
				len(message.Value), message.Timestamp, message.Topic, message.Partition, message.Offset, message.Headers)
			for _, b := range message.Value {
				consumer.Msgs <- b
				log.Printf("ConsumeClaim: value(%v)\n", b)
				if session.Context().Err() != nil {
					log.Printf("ConsumeClaim: session context error: %v", session.Context().Err())
					return nil
				}
			}
			if len(message.Value) == 0 {
				log.Printf("ConsumeClaim: message value is empty")
			}
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			log.Printf("ConsumeClaim: session context done err: %v\n", session.Context().Err())
			return nil
			// case <-consumer.doneCH:
			// 	log.Printf("ConsumeClaim: done channel done\n")
			// 	return nil
		}
	}
}
