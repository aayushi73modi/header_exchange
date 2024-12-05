package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare a header exchange
	err = ch.ExchangeDeclare(
		"header_exchange", // name
		"headers",         // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %s", err)
	}

	// Publish "type 1" message
	headers1 := amqp.Table{"type": "1"}
	body1 := "Message of type 1"
	err = ch.Publish(
		"header_exchange", // exchange
		"",                // routing key (ignored for headers)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body1),
			Headers:     headers1,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish type 1 message: %s", err)
	}
	log.Printf(" [x] Sent: %s", body1)

	// Publish "type 2" message
	headers2 := amqp.Table{"type": "2"}
	body2 := "Message of type 2"
	err = ch.Publish(
		"header_exchange", // exchange
		"",                // routing key (ignored for headers)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body2),
			Headers:     headers2,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish type 2 message: %s", err)
	}
	log.Printf(" [x] Sent: %s", body2)
}
