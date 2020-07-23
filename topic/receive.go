package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// ExchangeName - Declare and assign exchange name
var ExchangeName = "topic_wowza"

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

		err = ch.ExchangeDeclare(
			ExchangeName, // name
			"topic",      // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		)
		failOnError(err, "Failed to declare an exchange")

		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		failOnError(err, "Failed to declare a queue")

		// if len(os.Args) < 2 {
		// 	log.Printf("Usage: %s [binding_key]...", os.Args[0])
		// 	os.Exit(0)
		// }

		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, ExchangeName, "topic_wowza.abc.*")

		err = ch.QueueBind(
			q.Name,              // queue name
			"topic_wowza.abc.*", // routing key
			ExchangeName,        // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto ack
			false,  // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Printf(" [x] %s", d.Body)
				time.Sleep(time.Second * 2)
				log.Printf(" [x] %s success", d.Body)
			}
		}()

		log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
		<-forever

	}()

	<-forever2
}
