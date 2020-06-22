package main

import (
	"bytes"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	forever2 := make(chan bool)

	go func() {

		conn, err := amqp.Dial("amqp://user:VimRkNp1VjGJ@35.186.158.107:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		// err = ch.ExchangeDeclare(
		// 	"testja2", // name
		// 	"direct",  // type
		// 	true,      // durable
		// 	false,     // auto-deleted
		// 	false,     // internal
		// 	false,     // no-wait
		// 	nil,       // arguments
		// )
		// failOnError(err, "Failed to declare a exchange")

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // args
		)
		failOnError(err, "Failed to declare a queue")

		// err = ch.QueueBind(
		// 	q.Name, // queue name
		// 	"",     // routing key
		// 	"logs", // exchange
		// 	false,
		// 	nil,
		// )

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				dotCount := bytes.Count(d.Body, []byte("."))
				t := time.Duration(dotCount)
				time.Sleep(t * time.Second)
				log.Printf("Done")
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever

	}()

	<-forever2
}
