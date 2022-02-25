package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting Client")

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
		"amq.rabbitmq.reply-to",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		go func() {
			id := getRandomID()
			err = ch.Publish(
				"",
				"foo",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(id),
					ReplyTo:     "amq.rabbitmq.reply-to",
				},
			)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Published message to queue %s", id)

			//wait := make(chan bool)
			//go func() {
			//	for msg := range msgs {
			//		log.Printf("Received reply: %s", msg.Body)
			//		wait <- false
			//	}
			//}()

			msg := <-msgs
			log.Printf("Received reply: %s", msg.Body)
		}()
	}
}

func getRandomID() string {
	return strconv.Itoa(rand.Intn(9999999999))
}
