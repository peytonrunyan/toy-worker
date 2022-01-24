package repository

import (
	"github.com/streadway/amqp"
)

// TODO
type Repository interface {
	Write([]byte) error
}

type RabbitMQ struct {
	connString  string
	conn        *amqp.Connection
	amqpChannel *amqp.Channel
	queue       amqp.Queue
	MessageChan <-chan amqp.Delivery
}

func (r *RabbitMQ) Write(body []byte) error {
	err := r.amqpChannel.Publish("", r.queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) Setup() error {
	conn, err := amqp.Dial(r.connString)
	if err != nil {
		return err
	}
	r.conn = conn

	// Establish a channel to work with RabbitMQ
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}
	r.amqpChannel = channel

	// Specify the queue within RMQ that we're interested in
	queue, err := r.amqpChannel.QueueDeclare(
		"testMessages", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	r.queue = queue

	// Setting to deliver new messages to an RMQ node only when we have
	// an acknowledgement that it has received the previous one.
	err = r.amqpChannel.Qos(1, 0, false)
	if err != nil {
		return err
	}

	// Generate go channel to communicate with the queue
	messageChan, err := r.amqpChannel.Consume(
		r.queue.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	r.MessageChan = messageChan

	return nil
}

// Handle closing all open channels
func (r *RabbitMQ) CleanUp() {
	r.conn.Close()
	r.amqpChannel.Close()
}

// Generate new RabbitMQ with connection and channel initialized
func NewRMQ(s string) (*RabbitMQ, error) {
	rmq := RabbitMQ{
		connString: s,
	}

	err := rmq.Setup()
	if err != nil {
		return nil, err
	}

	return &rmq, nil
}
