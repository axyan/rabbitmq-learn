package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting Server")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	log.Println("Connected to RabbitMQ instance")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"foo",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	wait := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", msg.Body)

			err = ch.Publish(
				"",
				msg.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("Hello back %s", msg.Body)),
				},
			)
			if err != nil {
				log.Fatalln(err)
			}

			if err := ch.Ack(msg.DeliveryTag, false); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	log.Println("Consuming message from queue")

	<-wait
}
