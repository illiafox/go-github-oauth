package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Session struct {
	db *pgxpool.Pool
}

func (s Session) New(token string, user int64) error {

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO sessions (token,user_id,created) VALUES ($1,$2,NOW())", token, user)

	return err
}

// Exists returns user id if session exists, otherwise -1
func (s Session) Exists(token string) (int64, error) {

	var id int64 = -1
	err := s.db.QueryRow(context.Background(), "SELECT user_id FROM sessions WHERE token=$1", token).Scan(&id)
	if err != nil && errors.As(err, &sql.ErrNoRows) {
		err = nil
	}

	return id, err
}

func (s Session) Delete(token string) error {
	_, err := s.db.Exec(context.Background(), "DELETE FROM sessions WHERE token=$1", token)

	return err
}

func (s Session) DeleteOld(days int) (int64, error) {
	tag, err := s.db.Exec(
		context.Background(),
		fmt.Sprintf("DELETE FROM sessions WHERE created < now() - '%d days' :: interval", days),
	)

	return tag.RowsAffected(), err
}
