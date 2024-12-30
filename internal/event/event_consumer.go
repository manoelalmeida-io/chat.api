package event

import (
	"chat_api/internal/model"
	"chat_api/internal/repository"
	"chat_api/internal/utils"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type EventConsumer struct {
	chatRepository        *repository.ChatRepository
	chatMessageRepository *repository.ChatMessageRepository
	broadcast             *chan model.ChatMessage
}

func NewEventConsumer(
	chatRepository *repository.ChatRepository,
	chatMessageRepository *repository.ChatMessageRepository,
	broadcast *chan model.ChatMessage,
) *EventConsumer {
	return &EventConsumer{chatRepository, chatMessageRepository, broadcast}
}

func (e *EventConsumer) ReceiveMessageSent(d amqp091.Delivery) {

	var command model.SendMessageCommand

	json.Unmarshal(d.Body, &command)

	log.Printf("Received a message: %s", d.Body)

	e.registerMessageSent(command)
	e.registerMessageReceived(command)
}

func (e *EventConsumer) registerMessageSent(command model.SendMessageCommand) error {
	chat, err := e.chatRepository.FindByUserRefAndUserId(command.To, command.FromUserId)
	if err != nil {
		return err
	}

	var chatId string

	if chat == nil {
		id, err := utils.GetSnowflakeInstance().GenerateId()
		if err != nil {
			return err
		}

		newChat := model.Chat{Id: id, UserRef: command.To, UserId: command.FromUserId}
		e.chatRepository.Save(newChat)
		chatId = newChat.Id
	} else {
		chatId = chat.Id
	}

	id, err := utils.GetSnowflakeInstance().GenerateId()
	if err != nil {
		return err
	}

	message := &model.ChatMessage{Id: id, Content: command.Message, UserRef: command.To, DeliveryType: "SENT", ChatId: chatId}
	e.chatMessageRepository.Save(message)

	return nil
}

func (e *EventConsumer) registerMessageReceived(command model.SendMessageCommand) error {
	chat, err := e.chatRepository.FindByUserRefAndUserId(command.From, command.ToUserId)
	if err != nil {
		return err
	}

	var chatId string

	if chat == nil {
		id, err := utils.GetSnowflakeInstance().GenerateId()
		if err != nil {
			return err
		}

		newChat := model.Chat{Id: id, UserRef: command.From, UserId: command.ToUserId}
		e.chatRepository.Save(newChat)
		chatId = newChat.Id
	} else {
		chatId = chat.Id
	}

	id, err := utils.GetSnowflakeInstance().GenerateId()
	if err != nil {
		return err
	}

	message := &model.ChatMessage{Id: id, Content: command.Message, UserRef: command.From, DeliveryType: "RECEIVED", ChatId: chatId}
	e.chatMessageRepository.Save(message)

	*e.broadcast <- *message

	return nil
}
