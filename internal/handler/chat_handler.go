package handler

import (
	"chat_api/internal/event"
	"chat_api/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	eventPublisher *event.EventPublisher
}

func NewChatHandler(eventPublisher *event.EventPublisher) *ChatHandler {
	return &ChatHandler{eventPublisher}
}

func (h *ChatHandler) SendMessageHandler(c echo.Context) error {
	userInfo := c.Get("userInfo").(*model.User)

	request := new(model.SendMessageCommandRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	command := new(model.SendMessageCommand)
	command.Message = request.Message
	command.From = userInfo.Email
	command.To = request.To
	command.FromUserId = userInfo.Id

	h.eventPublisher.SendMessage(*command)

	return c.NoContent(http.StatusOK)
}
