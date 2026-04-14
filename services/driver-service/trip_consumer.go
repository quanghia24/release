package main

import (
	"context"
	"log"
	"ride-sharing/shared/messaging"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type tripConsumer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ) *tripConsumer {
	return &tripConsumer{
		rabbitmq: rabbitmq,
	}
}

func (c *tripConsumer) Listen() error {
	ctx := context.Background()
	return c.rabbitmq.ConsumeMessages(ctx, c.rabbitmq.Queue.Name, func(ctx context.Context, msg amqp.Delivery) error {
		log.Printf("Driver received message: %v", string(msg.Body))
		time.Sleep(5 * time.Second)
		log.Printf("Driver done prossing message: %v", msg.Body)

		return nil
	})
}
