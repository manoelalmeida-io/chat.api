package repository

import (
	"chat_api/internal/model"
	"database/sql"
)

type UserContactRepository struct {
	db *sql.DB
}

func NewUserContactRepository(db *sql.DB) *UserContactRepository {
	return &UserContactRepository{db: db}
}

func (r *UserContactRepository) FindById(id int64) (*model.UserContact, error) {
	row := r.db.QueryRow("SELECT * FROM user_contact WHERE id = ?", id)

	var userContact model.UserContact

	if err := row.Scan(&userContact.Id, &userContact.FirstName, &userContact.LastName, &userContact.Email, &userContact.UserId); err != nil {
		return nil, err
	}

	return &userContact, nil
}

func (r *UserContactRepository) FindByUserId(userId int64) ([]model.UserContact, error) {
	userContacts := make([]model.UserContact, 0)

	rows, err := r.db.Query("SELECT * FROM user_contact WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var userContact model.UserContact

		if err := rows.Scan(&userContact.Id, &userContact.FirstName, &userContact.LastName, &userContact.Email, &userContact.UserId); err != nil {
			return nil, err
		}

		userContacts = append(userContacts, userContact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userContacts, nil
}

func (r *UserContactRepository) Save(userContact model.UserContact) (*model.UserContact, error) {
	var res sql.Result
	var err error

	if userContact.Id != 0 {
		res, err = r.db.Exec(
			"UPDATE user_contact SET first_name = ?, last_name = ?, email = ?, user_id = ? WHERE id = ?",
			userContact.FirstName, userContact.LastName, userContact.Email, userContact.UserId, userContact.Id)
	} else {
		res, err = r.db.Exec(
			"INSERT INTO user_contact (first_name, last_name, email, user_id) VALUES (?, ?, ?, ?)",
			userContact.FirstName, userContact.LastName, userContact.Email, userContact.UserId)
	}

	if err != nil {
		return nil, err
	}

	if userContact.Id != 0 {
		return r.FindById(userContact.Id)
	}

	lastInsertedId, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	return r.FindById(lastInsertedId)
}

func (r *UserContactRepository) DeleteById(id int64) error {
	_, err := r.db.Exec("DELETE FROM user_contact WHERE id = ?", id)

	return err
}
