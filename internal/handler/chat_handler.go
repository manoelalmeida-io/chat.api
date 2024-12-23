package handler

import (
	"chat_api/internal/event"
	"chat_api/internal/model"
	"chat_api/internal/repository"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	eventPublisher *event.EventPublisher
	userRepository *repository.UserRepository
}

func NewChatHandler(eventPublisher *event.EventPublisher, userRepository *repository.UserRepository) *ChatHandler {
	return &ChatHandler{eventPublisher, userRepository}
}

func (h *ChatHandler) SendMessageHandler(c echo.Context) error {
	userInfo := c.Get("userInfo").(*model.User)

	request := new(model.SendMessageCommandRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	toUser, err := h.userRepository.FindByEmail(request.To)
	if err != nil && err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusBadRequest, "receiver user not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	command := new(model.SendMessageCommand)
	command.Message = request.Message
	command.From = userInfo.Email
	command.To = request.To
	command.FromUserId = userInfo.Id
	command.ToUserId = toUser.Id

	h.eventPublisher.SendMessage(*command)

	return c.NoContent(http.StatusOK)
}
