package db

import (
	"fmt"
	"server/config"
	"server/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection

func connectRabbitMQ(cfg *config.Config) (*amqp.Connection, error) {
	dsn := cfg.GetRabbitMQDSN()
	utils.Logger.Infof("use RabbitMQ DSN:%s", dsn)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// GetChannel creates and returns a new channel from the connection
func GetChannel() (*amqp.Channel, error) {
	if RabbitMQConn == nil {
		return nil, fmt.Errorf("RabbitMQ connection is not initialized")
	}
	return RabbitMQConn.Channel()
}

// WorkQueue represents a work queue pattern implementation
type WorkQueue struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// ========== Work Queue Pattern ==========
// NewWorkQueue creates a new work queue
func NewWorkQueue(queueName string) (*WorkQueue, error) {
	ch, err := GetChannel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	return &WorkQueue{
		Channel: ch,
		Queue:   q,
	}, nil
}

// Publish publishes a message to the work queue
func (wq *WorkQueue) Publish(message []byte) error {
	return wq.Channel.Publish(
		"",            // exchange
		wq.Queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		})
}

// Consume starts consuming messages from the work queue
func (wq *WorkQueue) Consume() (<-chan amqp.Delivery, error) {
	return wq.Channel.Consume(
		wq.Queue.Name, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
}

// Close closes the work queue channel
func (wq *WorkQueue) Close() error {
	return wq.Channel.Close()
}

// PublishSubscribe represents a publish/subscribe pattern implementation
type PublishSubscribe struct {
	Channel  *amqp.Channel
	Exchange string
}

// ========== Publish/Subscribe Pattern ==========
//
// NewPublishSubscribe creates a new publish/subscribe setup
func NewPublishSubscribe(exchangeName string) (*PublishSubscribe, error) {
	ch, err := GetChannel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	return &PublishSubscribe{
		Channel:  ch,
		Exchange: exchangeName,
	}, nil
}

// Publish broadcasts a message to all bound queues
func (ps *PublishSubscribe) Publish(message []byte) error {
	return ps.Channel.Publish(
		ps.Exchange, // exchange
		"",          // routing key (ignored for fanout)
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		})
}

// Subscribe creates a temporary queue and binds it to the exchange
func (ps *PublishSubscribe) Subscribe() (<-chan amqp.Delivery, error) {
	// Create an exclusive queue with a random name
	q, err := ps.Channel.QueueDeclare(
		"",    // name (empty = auto-generated)
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	// Bind the queue to the exchange
	err = ps.Channel.QueueBind(
		q.Name,      // queue name
		"",          // routing key (ignored for fanout)
		ps.Exchange, // exchange
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
	}

	return ps.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

// Close closes the publish/subscribe channel
func (ps *PublishSubscribe) Close() error {
	return ps.Channel.Close()
}

// Routing represents a routing pattern implementation
type Routing struct {
	Channel  *amqp.Channel
	Exchange string
}

// ========== Routing Pattern ==========
//
// NewRouting creates a new routing setup
func NewRouting(exchangeName string) (*Routing, error) {
	ch, err := GetChannel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	return &Routing{
		Channel:  ch,
		Exchange: exchangeName,
	}, nil
}

// PublishWithKey publishes a message with a specific routing key
func (r *Routing) PublishWithKey(routingKey string, message []byte) error {
	return r.Channel.Publish(
		r.Exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		})
}

// SubscribeWithKey subscribes to messages with a specific routing key
func (r *Routing) SubscribeWithKey(routingKey, queueName string) (<-chan amqp.Delivery, error) {
	q, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	err = r.Channel.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		r.Exchange, // exchange
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}

	return r.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

// Close closes the routing channel
func (r *Routing) Close() error {
	return r.Channel.Close()
}

// Topic represents a topic pattern implementation
type Topic struct {
	Channel  *amqp.Channel
	Exchange string
}

// ========== Topic Pattern ==========

// NewTopic creates a new topic setup
func NewTopic(exchangeName string) (*Topic, error) {
	ch, err := GetChannel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	return &Topic{
		Channel:  ch,
		Exchange: exchangeName,
	}, nil
}

// PublishWithPattern publishes a message with a topic pattern
func (t *Topic) PublishWithPattern(pattern string, message []byte) error {
	return t.Channel.Publish(
		t.Exchange, // exchange
		pattern,    // routing key (topic pattern)
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		})
}

// SubscribeWithPattern subscribes to messages matching a topic pattern
func (t *Topic) SubscribeWithPattern(pattern, queueName string) (<-chan amqp.Delivery, error) {
	q, err := t.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	err = t.Channel.QueueBind(
		q.Name,     // queue name
		pattern,    // routing key (topic pattern)
		t.Exchange, // exchange
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}

	return t.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

// Close closes the topic channel
func (t *Topic) Close() error {
	return t.Channel.Close()
}
