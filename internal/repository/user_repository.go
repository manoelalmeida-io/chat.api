package repository

import (
	"chat_api/internal/model"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindById(id int64) (*model.User, error) {
	row := r.db.QueryRow("SELECT * FROM user WHERE id = ?", id)

	var user model.User

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.GoogleSub); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindBySub(sub string) (*model.User, error) {
	row := r.db.QueryRow("SELECT * FROM user WHERE google_sub = ?", sub)

	var user model.User

	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.GoogleSub)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Save(user model.User) (*model.User, error) {
	var res sql.Result
	var err error

	if user.Id != 0 {
		res, err = r.db.Exec(
			"UPDATE user SET first_name = ?, last_name = ?, email = ?, google_sub = ? WHERE id = ?",
			user.FirstName, user.LastName, user.Email, user.GoogleSub, user.Id)
	} else {
		res, err = r.db.Exec(
			"INSERT INTO user (first_name, last_name, email, google_sub) VALUES (?, ?, ?, ?)",
			user.FirstName, user.LastName, user.Email, user.GoogleSub)
	}

	if err != nil {
		return nil, err
	}

	if user.Id != 0 {
		return r.FindById(user.Id)
	}

	lastInsertedId, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	return r.FindById(lastInsertedId)
}
