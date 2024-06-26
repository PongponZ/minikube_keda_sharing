package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	godotenv.Load(".env")

	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	publishQueue := os.Getenv("PUBLISH_QUEUE")
	sleepTimeEnv := os.Getenv("SLEEP_TIME")
	sleepTime, err := strconv.Atoi(sleepTimeEnv)
	if err != nil {
		log.Fatalf("Failed to convert sleep time to integer: %s", err)
	}

	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort)
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	ch.Qos(100, 0, false)

	messages, err := ch.Consume(
		publishQueue, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	blockProcess := make(chan bool)
	bulk := []amqp.Delivery{}
	limit := 100

	waitTimeOut := time.After(10 * time.Second)

	for message := range messages {
		select {
		case <-waitTimeOut:
			if len(bulk) > 0 {
				for _, m := range bulk {
					m.Ack(false)
				}
			}
		default:
			bulk = append(bulk, message)
			if len(bulk) == limit {
				for _, message := range bulk {
					message.Ack(false)
				}
				bulk = []amqp.Delivery{}
				time.Sleep(time.Duration(sleepTime) * time.Second)
				waitTimeOut = time.After(10 * time.Second)
			}
		}
	}

	<-blockProcess
}
