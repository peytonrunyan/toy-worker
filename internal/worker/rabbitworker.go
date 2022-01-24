package worker

import (
	c "dataprocessor/internal/common"
	r "dataprocessor/internal/repository"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type rabbitWorker struct {
	repo *r.RabbitMQ
}

// Processes messages received from rabbitMQ into RawMessages
func (w *rabbitWorker) process(msg amqp.Delivery) (*c.RawMessage, error) {
	rawMsg := &c.RawMessage{}
	if err := json.Unmarshal(msg.Body, rawMsg); err != nil {
		return nil, err
	}
	if err := msg.Ack(false); err != nil {
		return nil, err
	}
	return rawMsg, nil
}

func (w *rabbitWorker) send(msg *c.RawMessage) {
	fmt.Printf("msg: %v\n", msg)
}

func (w *rabbitWorker) Consume(wg *sync.WaitGroup) {
	c := make(chan *c.RawMessage)

	go func() {
		for val := range w.repo.MessageChan {
			msg, _ := w.process(val)
			c <- msg
		}
	}()

	go func() {
		for processedVal := range c {
			w.send(processedVal)
		}
		wg.Done()
	}()
}

func (w *rabbitWorker) Produce(wg *sync.WaitGroup) {
	go func() {
		i := 0
		for {
			rawMsg := &c.RawMessage{
				Msg: "Woof",
				ID:  i,
			}
			body, _ := json.Marshal(rawMsg)
			w.repo.Write(body)
			fmt.Printf("Produced message %d\n", i)
			i++
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func NewRabbitWorker(r *r.RabbitMQ) *rabbitWorker {
	w := rabbitWorker{
		repo: r,
	}
	return &w
}
