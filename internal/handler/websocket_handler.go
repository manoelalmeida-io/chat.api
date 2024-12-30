package handler

import (
	"chat_api/internal/model"
	"log"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type WebSocketHandler struct {
	clients   map[*websocket.Conn]bool
	Broadcast chan model.ChatMessage
}

func NewWebSocketHandler(broadcast chan model.ChatMessage) *WebSocketHandler {
	return &WebSocketHandler{
		clients:   make(map[*websocket.Conn]bool),
		Broadcast: broadcast,
	}
}

func (h *WebSocketHandler) Connect(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		h.clients[ws] = true

		for {
			var msg model.ChatMessage
			err := websocket.JSON.Receive(ws, &msg)
			if err != nil {
				delete(h.clients, ws)
				c.Logger().Error(err)
				return
			}
		}

	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

func (h *WebSocketHandler) HandleMessages() {
	for {
		msg := <-h.Broadcast
		for client := range h.clients {
			err := websocket.JSON.Send(client, msg)
			if err != nil {
				log.Printf("Error sending websocket message: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}
