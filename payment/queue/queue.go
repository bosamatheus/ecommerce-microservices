package queue

import (
	"fmt"

	amqp "github.com/streadway/amqp"
)

func failOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Connect() *amqp.Channel {
	// dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")
	dsn := "amqp://rabbitmq:rabbitmq@localhost:5672/"

	conn, err := amqp.Dial(dsn)
	failOnError(err)

	ch, err := conn.Channel()
	failOnError(err)

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

	failOnError(err)

	fmt.Println("Message sent")
}

func StartConsuming(queue string, ch *amqp.Channel, in chan []byte) {

	q, err := ch.QueueDeclare(queue,
		true,
		false,
		false,
		false,
		nil)

	msgs, err := ch.Consume(q.Name,
		"checkout",
		true,
		false,
		false,
		false,
		nil)

	failOnError(err)

	// Filling the channel in
	go func() {
		for m := range msgs {
			in <- []byte(m.Body)
		}
		close(in)
	}()
}
