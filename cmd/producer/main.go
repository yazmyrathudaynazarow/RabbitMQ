package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"programmingpercy/eventdrivenrmq/internal"
	"time"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("percy", "secret", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	if err := client.CreateQueue("customers_created", true, false); err != nil {
		panic(err)
	}

	if err := client.CreateQueue("customers_test", true, false); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_created", "customers.created.*", "customer_events"); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_test", "customer.*", "customer_events"); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Send(ctx, "customer_events", "customers.created.us", amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(`A cool message between servers 1`),
	}); err != nil {
		panic(err)
	}

	if err := client.Send(ctx, "customer_events", "customer.salam", amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(`An uncool message between servers`),
	}); err != nil {
		panic(err)
	}

	log.Println(client)
}
