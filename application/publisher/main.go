package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	publishQueue := os.Getenv("PUBLISH_QUEUE")

	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort)
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	queue, err := ch.QueueDeclare(
		publishQueue, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!, I am a publisher!")
	})

	app.Get("/publish/:quantity", func(c *fiber.Ctx) error {
		quantityString := c.Params("quantity")

		quantity, err := strconv.Atoi(quantityString)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid quantity",
			})
		}

		for i := 0; i < quantity; i++ {
			text := "Hello, World!, I am a publisher!"
			err = ch.Publish(
				"",         // exchange
				queue.Name, // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(text),
				})
			if err != nil {
				log.Printf("Failed to publish a message: %s\n", err)
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "Published!",
			"quantity": quantity,
		})
	})

	app.Listen(port)
}
