package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		panic(err.Error())
	}
}

func Connect() *amqp.Channel {
	// dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")
	dsn := "amqp://rabbitmq:rabbitmq@localhost:5672/"

	conn, err := amqp.Dial(dsn)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return ch
}

func Notify(payload []byte, exchange string, routingKey string, ch *amqp.Channel) {
	err := ch.Publish(exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		})

	failOnError(err, "Failed to notify")

	fmt.Println("Message sent")
}
