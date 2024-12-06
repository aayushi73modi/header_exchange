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

	// Consume messages from the queue
	queueName := "queue101"
	msgs, err := ch.Consume(
		queueName, // queue name
		"",        // consumer name
		true,      // auto-ack (set to false so we can ack manually)
		false,     // exclusive (this consumer is not exclusive)
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
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
			//msg.Ack(true) // Acknowledge the message
		} else {
			// Reject and requeue if the message is not of type 2
			//msg.Nack(false, true)
		}
	}
}
