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

func (r *ChatRepository) FindByUserId(userId int64) ([]model.Chat, error) {
	chats := make([]model.Chat, 0)

	rows, err := r.db.Query("SELECT c.*, uc.* FROM chat c LEFT JOIN user_contact uc ON c.user_ref = uc.email AND uc.user_id = ? WHERE c.user_id = ?", userId, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var chat model.Chat
		var contact model.UserContact

		if err := rows.Scan(&chat.Id, &chat.UserRef, &chat.UserId, &contact.Id,
			&contact.FirstName, &contact.LastName, &contact.Email, &contact.UserId); err != nil {
			return nil, err
		}

		if contact.Id != nil {
			chat.Contact = &contact
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
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
