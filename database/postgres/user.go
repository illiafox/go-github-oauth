package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	db *pgxpool.Pool
}

// UpdateUsername set new username for user
func (u User) UpdateUsername(id int64, username string) error {

	_, err := u.db.Exec(context.Background(), "UPDATE users SET username = $1 WHERE user_id = $2", username, id)

	return err
}

// Exists returns user id if it exists, otherwise -1
func (u User) Exists(id int64) (string, error) {

	var username string
	err := u.db.QueryRow(context.Background(), "SELECT username FROM users WHERE user_id=$1", id).Scan(&username)
	if err != nil && errors.As(err, &sql.ErrNoRows) {
		err = nil
	}

	return username, err
}

// Create inserts new user, returns only internal error
func (u User) Create(id int64, token, login string) error {

	_, err := u.db.Exec(context.Background(), "INSERT INTO users VALUES ($1,$2,$3)", id, login, token)
	return err
}

func (u User) Username(id int64) (string, error) {
	var username string

	err := u.db.QueryRow(context.Background(),
		"SELECT username FROM users WHERE user_id = $1", id).
		Scan(&username)

	return username, err
}
