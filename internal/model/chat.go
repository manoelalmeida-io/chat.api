package model

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
