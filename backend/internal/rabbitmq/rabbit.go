package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	fileUploadExchange string = "fileUploadExchange"
	AIContentGenQueue  string = "aiContentGen"
)

type RabbitMQ struct {
	conn *amqp.Connection
}

func New(user, password, host string) *RabbitMQ {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:5672/", user, password, host)
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("Failed to connect RabbitMQ")
	}
	m := &RabbitMQ{conn: conn}
	m.setup()
	return m
}

// here should create all the exchanges and the queue needed
func (rbmq *RabbitMQ) setup() {
	//create all the exchanges
	ch, err := rbmq.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	ch.ExchangeDeclare(
		fileUploadExchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	ch.QueueDeclare(AIContentGenQueue,
		true,
		false,
		false,
		false,
		nil,
	)

	ch.QueueBind(AIContentGenQueue, "", fileUploadExchange, false, nil)

}

func (rbmq *RabbitMQ) Publish(ctx context.Context, objectID uuid.UUID) error {
	ch, err := rbmq.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to craete channel: %w", err)
	}
	err = ch.PublishWithContext(ctx,
		fileUploadExchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(objectID.String()),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)

	}
	return nil
}

// what should we do,if we do consume here, it will be too tight coupled,
func (rbmq *RabbitMQ) Consume() chan amqp.Delivery {
	ch, _ := rbmq.conn.Channel()
	outCh := make(chan amqp.Delivery, 20)

	go func() {
		msgs, _ := ch.Consume(AIContentGenQueue, "", false, false, false, false, nil)
		for msg := range msgs {
			//now workers will get the msg, and ack it by themselfves
			outCh <- msg
		}
	}()
	return outCh
}
