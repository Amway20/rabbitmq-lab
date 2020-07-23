package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	for i := 0; i < 5; i++ {

		conn, err := amqp.Dial("amqp://user:XB2BNSTrcM@34.87.139.45:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"wow_za", // name
			false,    // durable
			false,    // delete when unused
			false,    // exclusive
			false,    // no-wait
			amqp.Table{
				"x-single-active-consumer": true,
				// "x-message-ttl":            int32(0),
			}, // args
		)

		failOnError(err, "Failed to declare a queue")

		qr, errQr := ch.QueueDeclare(
			"wow_za_reply", // name
			false,          // durable
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			amqp.Table{
				"x-single-active-consumer": true,
				// "x-message-ttl":            int32(0),
			}, // args
		)

		failOnError(errQr, "Failed to declare a queue")

		body := bodyFrom(os.Args)
		body += fmt.Sprintf(" %d", i)

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})

		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)

		msgs, err := ch.Consume(
			qr.Name, // queue
			"",      // consumer
			true,    // auto-ack
			false,   // exclusive
			false,   // no-local
			false,   // no-wait
			nil,     // args
		)
		failOnError(err, "Failed to register a consumer")

		// go func() {
		for d := range msgs {
			ch.Close()
			log.Printf("Reply %s", string(d.Body))
		}
		// }()
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
