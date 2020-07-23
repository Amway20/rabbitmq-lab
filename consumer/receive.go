package main

import (
	"bytes"
	"fmt"
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
			false,  // auto-ack
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
				fmt.Println(string(d.Body))
				if string(d.Body) == " 3" {
					log.Println("eieiei")
					break
				} else {

				}
				dotCount := bytes.Count(d.Body, []byte("."))
				t := time.Duration(dotCount)
				time.Sleep(t * time.Second)
				log.Printf("Done")

				result := []byte(`{"success":true"}`)
				d.Ack(false)

				err = ch.Publish(
					"",      // exchange
					qr.Name, // routing key
					false,   // mandatory
					false,   // immediate
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
						Body:         result,
					})

				failOnError(err, "Failed to publish a message")
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever

	}()

	<-forever2
}
