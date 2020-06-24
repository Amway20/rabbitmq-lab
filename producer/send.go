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
	forever := make(chan bool)
	for i := 0; i < 5; i++ {

		conn, err := amqp.Dial("amqp://user:XB2BNSTrcM@34.87.139.45:5672/")
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
			"test_queue_create_facebook_live_order", // name
			false,                                   // durable
			false,                                   // delete when unused
			false,                                   // exclusive
			false,                                   // no-wait
			amqp.Table{
				"x-single-active-consumer": true,
			}, // args
		)

		failOnError(err, "Failed to declare a queue")

		qr, errQr := ch.QueueDeclare(
			"queue_create_facebook_live_order_reply", // name
			false,                                    // durable
			false,                                    // delete when unused
			false,                                    // exclusive
			false,                                    // no-wait
			amqp.Table{
				"x-single-active-consumer": true,
			}, // args
		)

		failOnError(errQr, "Failed to declare a queue")

		// err = ch.Qos(
		// 	1,     // prefetch count
		// 	0,     // prefetch size
		// 	false, // global
		// )
		// failOnError(err, "Failed to set QoS")

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

	<-forever
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
