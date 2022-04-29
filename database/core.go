package database

import (
	"fmt"

	"go.uber.org/zap"
	"oauth/database/memcached"
	"oauth/database/postgres"
	"oauth/oauth"
	"oauth/utils/config"
)

type PoolClose func()

type Database struct {
	Memcached *memcached.Memcached
	//
	Postgres *postgres.Postgres
	//
	Oauth *oauth.Oauth
	//
	closePSQL PoolClose
}

func (d Database) Close(logger *zap.Logger) {
	logger.Info("Closing database connections")

	d.closePSQL()
	//
}

func New(conf *config.Config) (*Database, error) {

	mc, err := memcached.New(conf.Memcached)
	if err != nil {
		return nil, fmt.Errorf("memcached: %w", err)
	}

	pg, closePG, err := postgres.New(conf.Postgres)
	if err != nil {
		return nil, fmt.Errorf("posgres: %w", err)
	}

	return &Database{
		Memcached: mc,
		Postgres:  pg,

		Oauth: oauth.New(conf.Oauth),

		closePSQL: closePG,
	}, nil

}
