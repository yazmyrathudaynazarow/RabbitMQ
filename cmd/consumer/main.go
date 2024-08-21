package main

import (
	"log"
	"rabbitmq-servers/internal"
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

	messageBus, err := client.Consume("customers_created", "email-service", false)

	var blocking chan struct{}

	go func() {
		for message := range messageBus {

			log.Println("New Message: %v", string(message.Body))
			message.Ack(false)

		}
	}()
	log.Println("Consuming, use CTRL + C to stop")
	<-blocking
}
