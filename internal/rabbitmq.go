package internal

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitClient struct {
	// The connection used by client
	conn *amqp.Connection
	// Channel is used to process / Send messages
	ch *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()

	if err != nil {
		return RabbitClient{}, err
	}

	if err := ch.Confirm(false); err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) CreateQueue(name string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(name,
		durable,
		autodelete,
		false,
		false,
		nil,
	)
	return err
}

func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	return rc.ch.QueueBind(name,
		binding,
		exchange,
		false,
		nil)
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx,
		exchange,
		routingKey,
		true,
		false,
		options,
	)
	if err != nil {
		return err
	}

	log.Println(confirmation.Wait())
	return nil
}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue,
		consumer,
		autoAck,
		false,
		false,
		false,
		nil)
}
