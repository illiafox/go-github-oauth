package postgres

import (
	"context"
	"fmt"
	"time"

	// postgres
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	//
	"oauth/utils/config"
)

type Postgres struct {
	User    *User
	Session *Session
}

func New(conf config.Postgres) (*Postgres, func(), error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, err := pgxpool.Connect(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%v:%v/%v?sslmode=disable",
			conf.User,
			conf.Pass,
			conf.IP,
			conf.Port,
			conf.DbName,
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("opening connection: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	return &Postgres{
		&User{pool},
		&Session{pool},
	}, pool.Close, nil
}
