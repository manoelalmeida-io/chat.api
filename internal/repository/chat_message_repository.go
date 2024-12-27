package repository

import (
	"chat_api/internal/model"
	"database/sql"
)

type ChatMessageRepository struct {
	db *sql.DB
}

func NewChatMessageRepository(db *sql.DB) *ChatMessageRepository {
	return &ChatMessageRepository{db}
}

func (r *ChatMessageRepository) FindById(id string) (*model.ChatMessage, error) {
	row := r.db.QueryRow("SELECT * FROM chat_message WHERE id = ?", id)

	var message model.ChatMessage

	if err := row.Scan(&message.Id, &message.Content, &message.UserRef, &message.DeliveryType, &message.ChatId); err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *ChatMessageRepository) FindByChatIdAndUserId(chatId string, userId int64) ([]model.ChatMessage, error) {
	messages := make([]model.ChatMessage, 0)

	rows, err := r.db.Query(
		"SELECT cm.* FROM chat_message cm INNER JOIN chat c ON cm.chat_id = c.id WHERE cm.chat_id = ? and c.user_id = ? ORDER BY id DESC", chatId, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var message model.ChatMessage

		if err := rows.Scan(&message.Id, &message.Content, &message.UserRef, &message.DeliveryType, &message.ChatId); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *ChatMessageRepository) Save(message *model.ChatMessage) (*model.ChatMessage, error) {
	var err error

	row := r.db.QueryRow("SELECT COUNT(*) FROM chat_message WHERE id = ?", message.Id)
	var count int64

	if err := row.Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		_, err = r.db.Exec(
			"UPDATE chat_message SET content = ?, user_ref = ?, delivery_type = ?, chat_id = ? WHERE id = ?",
			message.Content, message.UserRef, message.DeliveryType, message.ChatId, message.Id)
	} else {
		_, err = r.db.Exec(
			"INSERT INTO chat_message (id, content, user_ref, delivery_type, chat_id) VALUES (?, ?, ?, ?, ?)",
			message.Id, message.Content, message.UserRef, message.DeliveryType, message.ChatId)
	}

	if err != nil {
		return nil, err
	}

	return r.FindById(message.Id)
}
