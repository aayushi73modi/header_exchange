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

	// Consume messages from the queue
	msgs, err := ch.Consume(
		queue.Name, // queue name
		"",         // consumer name
		false,      // auto-ack (set to false so we can ack manually)
		false,      // exclusive (this consumer is not exclusive)
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	log.Println(" [*] Consumer 2 waiting for 'type=2' messages. To exit press CTRL+C")
	for msg := range msgs {
		// Check if the message has a "type=2" header
		msgType, ok := msg.Headers["type"].(string)
		if ok && msgType == "2" {
			log.Printf(" [x] Consumer 2 processed message: %s", msg.Body)
			msg.Ack(true) // Acknowledge the message
		} else {
			// Reject and requeue if the message is not of type 2
			msg.Nack(false, true)
		}
	}
}
