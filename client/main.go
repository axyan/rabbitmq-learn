package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	_, err = ch.QueueDeclare(
		"testQueue",
		false,
		false,
		false,
		false,
		nil,
	)

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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

requestsLoop:
	for {
		select {
		case <-stop:
			log.Println("Stopping client")
			break requestsLoop
		default:
			time.Sleep(getRandomSecondUpTo(5))
			go func() {
				id := getRandomID()
				err = ch.Publish(
					"",
					"testQueue",
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(fmt.Sprintf("Hello from %s", id)),
						ReplyTo:     "amq.rabbitmq.reply-to",
					},
				)
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("Published new message to queue for: %s", id)

				msg := <-msgs
				log.Printf("Received reply: %s", msg.Body)
			}()
		}
	}
}

func getRandomID() string {
	return strconv.Itoa(rand.Intn(9999999999))
}

func getRandomSecondUpTo(upperBound int) time.Duration {
	return time.Duration(rand.Intn(upperBound)) * time.Second
}
