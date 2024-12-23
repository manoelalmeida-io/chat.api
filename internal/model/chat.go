package model

type Chat struct {
	Id      string `json:"id"`
	UserRef string `json:"userRef"`
	UserId  int64  `json:"userId"`
}

type ChatMessage struct {
	Id           string `json:"id"`
	Content      string `json:"content"`
	UserRef      string `json:"userRef"`
	DeliveryType string `json:"deliveryType"`
	ChatId       string `json:"chatId"`
}

type SendMessageCommandRequest struct {
	Message string `json:"message"`
	To      string `json:"to"`
}

type SendMessageCommand struct {
	Message    string `json:"message"`
	From       string `json:"from"`
	To         string `json:"to"`
	FromUserId int64  `json:"fromUserId"`
	ToUserId   int64  `json:"toUserId"`
}
