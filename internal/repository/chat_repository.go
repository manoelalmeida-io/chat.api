package repository

import (
	"chat_api/internal/model"
	"database/sql"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db}
}

func (r *ChatRepository) FindById(id string) (*model.Chat, error) {
	row := r.db.QueryRow("SELECT * FROM chat WHERE id = ?", id)

	var chat model.Chat

	if err := row.Scan(&chat.Id, &chat.UserRef, &chat.UserId); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) FindByUserRefAndUserId(userRef string, userId int64) (*model.Chat, error) {
	row := r.db.QueryRow("SELECT * FROM chat WHERE user_ref = ? AND user_id = ?", userRef, userId)

	var chat model.Chat

	err := row.Scan(&chat.Id, &chat.UserRef, &chat.UserId)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) Save(chat model.Chat) (*model.Chat, error) {
	var err error

	row := r.db.QueryRow("SELECT COUNT(*) FROM chat WHERE id = ?", chat.Id)
	var count int64

	if err := row.Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		_, err = r.db.Exec(
			"UPDATE chat SET user_ref = ?, user_id = ? WHERE id = ?",
			chat.UserRef, chat.UserId, chat.Id)
	} else {
		_, err = r.db.Exec(
			"INSERT INTO chat (id, user_ref, user_id) VALUES (?, ?, ?)",
			chat.Id, chat.UserRef, chat.UserId)
	}

	if err != nil {
		return nil, err
	}

	return r.FindById(chat.Id)
}
