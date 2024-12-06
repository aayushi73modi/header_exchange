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
	queuName := "queue101"
	msgs, err := ch.Consume(
		queuName, // queue name
		"",       // consumer name
		true,     // auto-ack (set to false so we can ack manually)
		false,    // exclusive (this consumer is not exclusive)
		false,    // no-local
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	log.Println(" [*] Consumer 1 waiting for 'type=1' messages. To exit press CTRL+C")
	for msg := range msgs {
		msgType, ok := msg.Headers["type"].(string)
		if ok && msgType == "1" {
			log.Printf(" [x] Consumer 1 processed message: %s", msg.Body)
			//msg.Ack(true) // Acknowledge the message
		} else {
			log.Printf(" [!] Ignored message: %s", msg.Body)
			// Reject and requeue if the message is not of type 1
			//msg.Nack(false, true)
		}
	}
}
