package event

import (
	"chat_api/internal/model"
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher struct {
	amqpChannel *amqp.Channel
}

func NewEventPublisher(amqpChannel *amqp.Channel) *EventPublisher {
	return &EventPublisher{amqpChannel}
}

func (p *EventPublisher) SendMessage(command model.SendMessageCommand) {
	q, err := p.amqpChannel.QueueDeclare(
		"chat.message.sent.queue", // name
		false,                     // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		log.Fatal("failed to declare queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	commandJson, err := json.Marshal(command)
	if err != nil {
		log.Fatal("failed to marshal command to JSON", err)
	}

	err = p.amqpChannel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        commandJson,
		})
	if err != nil {
		log.Print("error publishing message to queue", err)
	}
}
