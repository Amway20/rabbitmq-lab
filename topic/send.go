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

// ExchangeName - Declare and assign exchange name
var ExchangeName = "topic_wowza"
var TopicName = "abc"

func main() {
	// forever := make(chan bool)
	for i := 0; i < 5; i++ {
		producerTopic(1)
	}

	// <-forever
}

func producerTopic(fbLiveID int) {
	fmt.Println("fbLiveID: ", fbLiveID)
	conn, err := amqp.Dial("amqp://user:XB2BNSTrcM@34.87.139.45:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	topic := fmt.Sprintf("%s.%s.%d", ExchangeName, TopicName, fbLiveID)

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

	// q, err := ch.QueueDeclare(
	// 	topic, // name
	// 	false, // durable
	// 	false, // delete when unused
	// 	true,  // exclusive
	// 	false, // no-wait
	// 	nil,   // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	// fmt.Println("q: ", q)

	// body := bodyFrom(os.Args)

	err = ch.Publish(
		ExchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%d", fbLiveID)),
		})

	failOnError(err, "Failed to publish a message")
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = ""
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
