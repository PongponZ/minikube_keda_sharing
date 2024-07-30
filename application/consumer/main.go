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
	waitTimeEnv := os.Getenv("WAIT_TIME")
	waitTime, err := strconv.Atoi(waitTimeEnv)
	if err != nil {
		log.Fatalf("Failed to convert wait time to integer: %s", err)
	}

	limitConsumeEnv := os.Getenv("LIMIT_CONSUME")
	limitConsume, err := strconv.Atoi(limitConsumeEnv)
	if err != nil {
		log.Fatalf("Failed to convert limit consume to integer: %s", err)
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

	bulk := []amqp.Delivery{}

	waitTimeOut := time.After(time.Duration(waitTime) * time.Second)

	for {
		select {
		case message := <-messages:
			bulk = append(bulk, message)
			if len(bulk) == limitConsume {
				for _, message := range bulk {
					message.Ack(false)
				}
				bulk = []amqp.Delivery{}
				time.Sleep(time.Duration(sleepTime) * time.Second)
				waitTimeOut = time.After(time.Duration(waitTime) * time.Second)
			}
		case <-waitTimeOut:
			if len(bulk) > 0 {
				for _, m := range bulk {
					m.Ack(false)
				}
			}
		}
	}
}
