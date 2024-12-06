package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
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

	// Declare the shared queue
	queue, err := ch.QueueDeclare(
		"queue101", // queue name
		true,       // durable
		false,      // auto-deleted
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	// Bind the queue to the header exchange for "type=1"
	err = ch.QueueBind(
		queue.Name,        // queue name
		"",                // routing key (ignored for headers)
		"header_exchange", // exchange
		false,             // no-wait
		amqp.Table{
			"x-match": "all", // Match all headers
			"type":    "1",   // Only match "type=1"
		},
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %s", err)
	}
	// Bind the queue to the header exchange for "type=2"
	err = ch.QueueBind(
		queue.Name,        // queue name
		"",                // routing key (ignored for headers)
		"header_exchange", // exchange
		false,             // no-wait
		amqp.Table{
			"x-match": "all", // Match all headers
			"type":    "2",   // Only match "type=2"
		},
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %s", err)
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
